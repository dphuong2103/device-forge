package agent

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	version  string
	sessions *SessionManager
	ui       fs.FS
}

func NewServer(version string, sessions *SessionManager, ui fs.FS) *Server {
	return &Server{version: version, sessions: sessions, ui: ui}
}

func (s *Server) Handler(ports []int) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", s.health)
	mux.HandleFunc("POST /api/v1/session/bootstrap", s.bootstrap)
	if s.ui != nil {
		mux.Handle("/", spaFileServer(s.ui))
	} else {
		mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "UI assets are not installed", http.StatusServiceUnavailable)
		})
	}
	security := newSecurityMiddleware(s.sessions, ports, 1<<20)
	return security.wrap(withSecurityHeaders(mux))
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok", "version": s.version})
}

func (s *Server) bootstrap(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Token string `json:"token"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil || strings.TrimSpace(input.Token) == "" {
		writeProblem(w, http.StatusBadRequest, "invalid-request", "Invalid request", "A bootstrap token is required.")
		return
	}
	sessionID, csrfToken, expiresAt, err := s.sessions.ExchangeBootstrap(input.Token)
	if err != nil {
		writeProblem(w, http.StatusUnauthorized, "invalid-bootstrap", "Invalid bootstrap token", "The token is invalid, expired, or already consumed.")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  expiresAt,
		MaxAge:   int(time.Until(expiresAt).Seconds()),
	})
	w.Header().Set(csrfHeaderName, csrfToken)
	w.WriteHeader(http.StatusNoContent)
}

func withSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; connect-src 'self' ws:; img-src 'self' data:; style-src 'self' 'unsafe-inline'; script-src 'self'; base-uri 'none'; frame-ancestors 'none'; form-action 'self'")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}

func spaFileServer(root fs.FS) http.Handler {
	files := http.FileServer(http.FS(root))
	index, indexErr := fs.ReadFile(root, "index.html")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}
		if info, err := fs.Stat(root, path); err == nil && !info.IsDir() {
			files.ServeHTTP(w, r)
			return
		}
		if indexErr != nil {
			http.Error(w, "UI assets are incomplete", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(index)
	})
}
