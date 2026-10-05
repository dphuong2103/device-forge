package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/device-force/device-force-v2/internal/agent"
	"github.com/device-force/device-force-v2/internal/observability"
)

var version = "dev"

//go:embed web/*
var webAssets embed.FS

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	defaultDataDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("resolve user config directory: %w", err)
	}
	port := flag.Int("port", 47631, "loopback API port")
	dataDir := flag.String("data-dir", filepath.Join(defaultDataDir, "DeviceForce"), "Agent data directory")
	flag.Parse()

	ui, err := fs.Sub(webAssets, "web")
	if err != nil {
		return fmt.Errorf("load UI assets: %w", err)
	}
	logger := observability.NewJSONLogger(os.Stderr, slog.LevelInfo)
	application, err := agent.New(agent.Config{
		Port:    *port,
		DataDir: *dataDir,
		Version: version,
		UI:      ui,
		Logger:  logger,
	})
	if err != nil {
		return err
	}
	bootstrap, err := application.Sessions().IssueBootstrap()
	if err != nil {
		return fmt.Errorf("create bootstrap token: %w", err)
	}
	logger.Info("Agent starting", "port", *port, "bootstrap_token", bootstrap)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err := application.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}
