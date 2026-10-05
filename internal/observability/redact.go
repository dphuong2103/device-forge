package observability

import (
	"log/slog"
	"strings"
)

const redacted = "[REDACTED]"

var sensitiveFragments = []string{
	"authorization",
	"cookie",
	"csrf",
	"password",
	"passwd",
	"proxy_url",
	"refresh_token",
	"secret",
	"session",
	"token",
}

func RedactAttr(_ []string, attr slog.Attr) slog.Attr {
	if attr.Equal(slog.Attr{}) {
		return attr
	}
	key := strings.ToLower(strings.ReplaceAll(attr.Key, "-", "_"))
	for _, fragment := range sensitiveFragments {
		if strings.Contains(key, fragment) {
			return slog.String(attr.Key, redacted)
		}
	}
	if attr.Value.Kind() == slog.KindGroup {
		group := attr.Value.Group()
		for i := range group {
			group[i] = RedactAttr(nil, group[i])
		}
		return slog.Group(attr.Key, attrsToAny(group)...)
	}
	return attr
}

func attrsToAny(attrs []slog.Attr) []any {
	values := make([]any, len(attrs))
	for i := range attrs {
		values[i] = attrs[i]
	}
	return values
}
