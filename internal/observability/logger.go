package observability

import (
	"context"
	"io"
	"log/slog"
)

func NewJSONLogger(output io.Writer, level slog.Leveler) *slog.Logger {
	return slog.New(slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level:       level,
		ReplaceAttr: RedactAttr,
	}))
}

func WithRequest(ctx context.Context, logger *slog.Logger) *slog.Logger {
	if id := CorrelationID(ctx); id != "" {
		return logger.With("correlation_id", id)
	}
	return logger
}

type Diagnostic struct {
	Code          string         `json:"code"`
	Severity      string         `json:"severity"`
	Message       string         `json:"message"`
	Remediation   string         `json:"remediation,omitempty"`
	CorrelationID string         `json:"correlationId"`
	Fields        map[string]any `json:"fields,omitempty"`
}
