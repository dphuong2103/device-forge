package agent

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
	"time"
)

func TestBootstrapEndpointIssuesProtectedSession(t *testing.T) {
	t.Parallel()
	sessions := NewSessionManager(time.Minute, time.Hour)
	token, _ := sessions.IssueBootstrap()
	handler := NewServer("test", sessions, testUI()).Handler([]int{47631})

	request := httptest.NewRequest(http.MethodPost, "http://localhost:47631/api/v1/session/bootstrap", strings.NewReader(`{"token":"`+token+`"}`))
	request.Host = "localhost:47631"
	request.Header.Set("Origin", "http://localhost:47631")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	cookies := response.Result().Cookies()
	if len(cookies) != 1 || !cookies[0].HttpOnly || cookies[0].SameSite != http.SameSiteStrictMode {
		t.Fatalf("unexpected session cookie: %#v", cookies)
	}
	if response.Header().Get(csrfHeaderName) == "" {
		t.Fatal("missing CSRF token")
	}

	replay := httptest.NewRequest(http.MethodPost, "http://localhost:47631/api/v1/session/bootstrap", strings.NewReader(`{"token":"`+token+`"}`))
	replay.Host = "localhost:47631"
	replay.Header.Set("Origin", "http://localhost:47631")
	replayResponse := httptest.NewRecorder()
	handler.ServeHTTP(replayResponse, replay)
	if replayResponse.Code != http.StatusUnauthorized {
		t.Fatalf("replay status = %d", replayResponse.Code)
	}
}

func TestSecurityRejectsHostOriginAndMissingCSRF(t *testing.T) {
	t.Parallel()
	sessions := NewSessionManager(time.Minute, time.Hour)
	handler := NewServer("test", sessions, testUI()).Handler([]int{47631})

	tests := []struct {
		name   string
		host   string
		origin string
		want   int
	}{
		{name: "DNS rebinding host", host: "evil.example:47631", origin: "http://evil.example:47631", want: http.StatusForbidden},
		{name: "hostile origin", host: "localhost:47631", origin: "https://evil.example", want: http.StatusForbidden},
		{name: "missing origin on mutation", host: "localhost:47631", origin: "", want: http.StatusForbidden},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "http://localhost:47631/api/v1/session/bootstrap", strings.NewReader(`{"token":"invalid"}`))
			request.Host = test.host
			if test.origin != "" {
				request.Header.Set("Origin", test.origin)
			}
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != test.want {
				t.Fatalf("status = %d, want %d", response.Code, test.want)
			}
		})
	}
}

func TestHealthAndStaticUI(t *testing.T) {
	t.Parallel()
	handler := NewServer("1.2.3", NewSessionManager(time.Minute, time.Hour), testUI()).Handler([]int{47631})

	health := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:47631/api/v1/health", nil)
	health.Host = "127.0.0.1:47631"
	healthResponse := httptest.NewRecorder()
	handler.ServeHTTP(healthResponse, health)
	if healthResponse.Code != http.StatusOK {
		t.Fatalf("health status = %d", healthResponse.Code)
	}
	var payload map[string]string
	if err := json.NewDecoder(healthResponse.Body).Decode(&payload); err != nil || payload["version"] != "1.2.3" {
		t.Fatalf("health payload = %#v, error = %v", payload, err)
	}

	ui := httptest.NewRequest(http.MethodGet, "http://[::1]:47631/devices/abc", nil)
	ui.Host = "[::1]:47631"
	uiResponse := httptest.NewRecorder()
	handler.ServeHTTP(uiResponse, ui)
	body, _ := io.ReadAll(uiResponse.Result().Body)
	if uiResponse.Code != http.StatusOK || !strings.Contains(string(body), "Device Force") {
		t.Fatalf("UI status = %d body = %s", uiResponse.Code, body)
	}
}

func testUI() fstest.MapFS {
	return fstest.MapFS{"index.html": {Data: []byte("<!doctype html><title>Device Force</title>")}}
}
