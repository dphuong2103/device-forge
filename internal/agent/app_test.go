package agent

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestAppBindsOnlyLoopbackAndStopsGracefully(t *testing.T) {
	port := freeDualStackPort(t)
	dataDir := t.TempDir()
	app, err := New(Config{Port: port, DataDir: dataDir, Version: "test", UI: testUI(), ShutdownTimeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- app.Run(ctx) }()

	waitForHealth(t, fmt.Sprintf("http://127.0.0.1:%d/api/v1/health", port))
	waitForHealth(t, fmt.Sprintf("http://[::1]:%d/api/v1/health", port))

	second, err := New(Config{Port: port, DataDir: dataDir, Version: "test", UI: testUI()})
	if err != nil {
		t.Fatal(err)
	}
	if err := second.Run(context.Background()); !errors.Is(err, ErrAlreadyRunning) {
		t.Fatalf("second instance error = %v", err)
	}

	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run returned %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Agent did not stop gracefully")
	}
}

func freeDualStackPort(t *testing.T) int {
	t.Helper()
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	_ = listener.Close()
	probe, err := net.Listen("tcp6", fmt.Sprintf("[::1]:%d", port))
	if err != nil {
		t.Skipf("IPv6 loopback unavailable: %v", err)
	}
	_ = probe.Close()
	return port
}

func waitForHealth(t *testing.T, endpoint string) {
	t.Helper()
	client := http.Client{Timeout: 200 * time.Millisecond}
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		response, err := client.Get(endpoint)
		if err == nil {
			_ = response.Body.Close()
			if response.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("Agent did not become healthy at %s", endpoint)
}
