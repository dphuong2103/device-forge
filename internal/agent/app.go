package agent

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gofrs/flock"
)

var ErrAlreadyRunning = errors.New("device force Agent is already running for this user")

type Config struct {
	Port            int
	DataDir         string
	Version         string
	UI              fs.FS
	Logger          *slog.Logger
	BootstrapTTL    time.Duration
	SessionTTL      time.Duration
	ShutdownTimeout time.Duration
}

type App struct {
	config   Config
	sessions *SessionManager
}

func New(config Config) (*App, error) {
	if config.Port < 1 || config.Port > 65535 {
		return nil, fmt.Errorf("invalid Agent port %d", config.Port)
	}
	if config.DataDir == "" {
		return nil, errors.New("data directory is required")
	}
	if config.Logger == nil {
		config.Logger = slog.New(slog.NewTextHandler(os.Stderr, nil))
	}
	if config.BootstrapTTL <= 0 {
		config.BootstrapTTL = 2 * time.Minute
	}
	if config.SessionTTL <= 0 {
		config.SessionTTL = 8 * time.Hour
	}
	if config.ShutdownTimeout <= 0 {
		config.ShutdownTimeout = 10 * time.Second
	}
	return &App{config: config, sessions: NewSessionManager(config.BootstrapTTL, config.SessionTTL)}, nil
}

func (a *App) Sessions() *SessionManager { return a.sessions }

func (a *App) Run(ctx context.Context) error {
	if err := os.MkdirAll(a.config.DataDir, 0o700); err != nil {
		return fmt.Errorf("create data directory: %w", err)
	}
	instanceLock := flock.New(filepath.Join(a.config.DataDir, "agent.lock"))
	locked, err := instanceLock.TryLock()
	if err != nil {
		return fmt.Errorf("acquire Agent lock: %w", err)
	}
	if !locked {
		return ErrAlreadyRunning
	}
	defer func() {
		if err := instanceLock.Unlock(); err != nil {
			a.config.Logger.Error("release Agent lock", "error", err)
		}
	}()

	addresses := []string{
		fmt.Sprintf("127.0.0.1:%d", a.config.Port),
		fmt.Sprintf("[::1]:%d", a.config.Port),
	}
	listeners := make([]net.Listener, 0, len(addresses))
	for _, address := range addresses {
		listener, err := net.Listen("tcp", address)
		if err != nil {
			for _, opened := range listeners {
				_ = opened.Close()
			}
			return fmt.Errorf("listen on %s: %w", address, err)
		}
		listeners = append(listeners, listener)
	}

	server := &http.Server{
		Handler:           NewServer(a.config.Version, a.sessions, a.config.UI).Handler([]int{a.config.Port}),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       90 * time.Second,
		MaxHeaderBytes:    32 << 10,
	}
	serveErrors := make(chan error, len(listeners))
	var serving sync.WaitGroup
	for _, listener := range listeners {
		serving.Add(1)
		go func() {
			defer serving.Done()
			if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
				serveErrors <- err
			}
		}()
	}

	var runErr error
	select {
	case <-ctx.Done():
	case runErr = <-serveErrors:
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.config.ShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil && runErr == nil {
		runErr = fmt.Errorf("shutdown Agent: %w", err)
	}
	serving.Wait()
	return runErr
}
