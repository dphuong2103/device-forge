package agent

import (
	"errors"
	"testing"
	"time"
)

func TestBootstrapCanBeExchangedOnlyOnce(t *testing.T) {
	t.Parallel()
	m := NewSessionManager(time.Minute, time.Hour)
	token, err := m.IssueBootstrap()
	if err != nil {
		t.Fatal(err)
	}
	session, csrf, _, err := m.ExchangeBootstrap(token)
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Validate(session, csrf, true); err != nil {
		t.Fatalf("valid session rejected: %v", err)
	}
	if _, _, _, err := m.ExchangeBootstrap(token); !errors.Is(err, ErrBootstrapInvalid) {
		t.Fatalf("expected replay rejection, got %v", err)
	}
}

func TestSessionRejectsWrongCSRFAndExpiry(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC)
	m := NewSessionManager(time.Minute, time.Hour)
	m.now = func() time.Time { return now }
	token, _ := m.IssueBootstrap()
	session, csrf, _, err := m.ExchangeBootstrap(token)
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Validate(session, "wrong", true); !errors.Is(err, ErrCSRFInvalid) {
		t.Fatalf("expected CSRF rejection, got %v", err)
	}
	if err := m.Validate(session, csrf, true); err != nil {
		t.Fatal(err)
	}
	now = now.Add(time.Hour + time.Nanosecond)
	if err := m.Validate(session, csrf, false); !errors.Is(err, ErrSessionInvalid) {
		t.Fatalf("expected expiry, got %v", err)
	}
}
