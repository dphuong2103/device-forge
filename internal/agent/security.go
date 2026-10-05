package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
)

const (
	sessionCookieName = "device_force_session"
	csrfHeaderName    = "X-CSRF-Token"
)

type securityMiddleware struct {
	sessions       *SessionManager
	allowedPorts   map[string]struct{}
	maxRequestBody int64
}

func newSecurityMiddleware(sessions *SessionManager, ports []int, maxRequestBody int64) *securityMiddleware {
	allowed := make(map[string]struct{}, len(ports))
	for _, port := range ports {
		allowed[fmt.Sprint(port)] = struct{}{}
	}
	return &securityMiddleware{sessions: sessions, allowedPorts: allowed, maxRequestBody: maxRequestBody}
}

func (m *securityMiddleware) wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.validHost(r.Host) {
			writeProblem(w, http.StatusForbidden, "untrusted-host", "Untrusted host", "The request host is not an Agent loopback address.")
			return
		}
		if !m.validOrigin(r) {
			writeProblem(w, http.StatusForbidden, "untrusted-origin", "Untrusted origin", "The request origin does not match the Agent origin.")
			return
		}
		if m.maxRequestBody > 0 && r.Body != nil {
			r.Body = http.MaxBytesReader(w, r.Body, m.maxRequestBody)
		}
		if isPublicPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil {
			writeProblem(w, http.StatusUnauthorized, "session-required", "Session required", "Open the interface from the Device Force tray application.")
			return
		}
		requireCSRF := isStateChanging(r.Method) || isWebSocketUpgrade(r)
		if err := m.sessions.Validate(cookie.Value, r.Header.Get(csrfHeaderName), requireCSRF); err != nil {
			status := http.StatusUnauthorized
			code := "invalid-session"
			if errors.Is(err, ErrCSRFInvalid) {
				status = http.StatusForbidden
				code = "invalid-csrf"
			}
			writeProblem(w, status, code, "Request authentication failed", err.Error())
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *securityMiddleware) validHost(value string) bool {
	host, port, err := net.SplitHostPort(value)
	if err != nil {
		return false
	}
	if _, ok := m.allowedPorts[port]; !ok {
		return false
	}
	return host == "127.0.0.1" || host == "::1" || strings.EqualFold(host, "localhost")
}

func (m *securityMiddleware) validOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return !isStateChanging(r.Method) && !isWebSocketUpgrade(r)
	}
	parsed, err := url.Parse(origin)
	if err != nil || parsed.Scheme != "http" || parsed.Host == "" {
		return false
	}
	return strings.EqualFold(parsed.Host, r.Host)
}

func isPublicPath(path string) bool {
	return path == "/api/v1/health" || path == "/api/v1/session/bootstrap" || !strings.HasPrefix(path, "/api/")
}

func isStateChanging(method string) bool {
	return method != http.MethodGet && method != http.MethodHead && method != http.MethodOptions
}

func isWebSocketUpgrade(r *http.Request) bool {
	return strings.EqualFold(r.Header.Get("Upgrade"), "websocket") && strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade")
}

type problem struct {
	Code   string `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail,omitempty"`
}

func writeProblem(w http.ResponseWriter, status int, code, title, detail string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(problem{Code: code, Title: title, Detail: detail})
}
