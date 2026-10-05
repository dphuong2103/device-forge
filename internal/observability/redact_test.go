package observability

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
)

func TestLoggerRedactsSensitiveAttributesRecursively(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	logger := NewJSONLogger(&output, slog.LevelDebug)
	logger.Info("proxy configured",
		"proxy_password", "hunter2",
		"device_id", "device-1",
		slog.Group("http", "Authorization", "Bearer abc", "status", 200),
	)

	line := output.String()
	for _, secret := range []string{"hunter2", "Bearer abc"} {
		if strings.Contains(line, secret) {
			t.Fatalf("log contains secret %q: %s", secret, line)
		}
	}
	if count := strings.Count(line, redacted); count != 2 {
		t.Fatalf("expected two redactions, got %d: %s", count, line)
	}
	if !strings.Contains(line, "device-1") || !strings.Contains(line, `"status":200`) {
		t.Fatalf("non-sensitive context was removed: %s", line)
	}
}

func TestCorrelationIDIsAddedToRequestLogger(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	logger := NewJSONLogger(&output, slog.LevelInfo)
	ctx := WithCorrelationID(context.Background(), "correlation-123456")
	WithRequest(ctx, logger).Info("request complete")

	if !strings.Contains(output.String(), `"correlation_id":"correlation-123456"`) {
		t.Fatalf("missing correlation ID: %s", output.String())
	}
}

func TestSensitiveKeyNormalization(t *testing.T) {
	t.Parallel()

	for _, key := range []string{"X-CSRF-Token", "refresh-token", "sessionCookie", "client_secret"} {
		attr := RedactAttr(nil, slog.String(key, "value"))
		if got := attr.Value.String(); got != redacted {
			t.Errorf("%s was not redacted: %q", key, got)
		}
	}
}
