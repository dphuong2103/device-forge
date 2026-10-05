package agent

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

var (
	ErrBootstrapInvalid = errors.New("bootstrap token is invalid, expired, or consumed")
	ErrSessionInvalid   = errors.New("session is invalid or expired")
	ErrCSRFInvalid      = errors.New("csrf token is invalid")
)

type clock func() time.Time

type bootstrapRecord struct {
	expiresAt time.Time
}

type sessionRecord struct {
	csrfHash  [sha256.Size]byte
	expiresAt time.Time
}

type SessionManager struct {
	mu           sync.Mutex
	now          clock
	bootstrapTTL time.Duration
	sessionTTL   time.Duration
	bootstraps   map[[sha256.Size]byte]bootstrapRecord
	sessions     map[[sha256.Size]byte]sessionRecord
}

func NewSessionManager(bootstrapTTL, sessionTTL time.Duration) *SessionManager {
	return &SessionManager{
		now:          time.Now,
		bootstrapTTL: bootstrapTTL,
		sessionTTL:   sessionTTL,
		bootstraps:   make(map[[sha256.Size]byte]bootstrapRecord),
		sessions:     make(map[[sha256.Size]byte]sessionRecord),
	}
}

func (m *SessionManager) IssueBootstrap() (string, error) {
	token, err := randomToken(32)
	if err != nil {
		return "", err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pruneLocked()
	m.bootstraps[hashToken(token)] = bootstrapRecord{expiresAt: m.now().Add(m.bootstrapTTL)}
	return token, nil
}

func (m *SessionManager) ExchangeBootstrap(token string) (sessionID, csrfToken string, expiresAt time.Time, err error) {
	hash := hashToken(token)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pruneLocked()
	record, ok := m.bootstraps[hash]
	if !ok || !record.expiresAt.After(m.now()) {
		delete(m.bootstraps, hash)
		return "", "", time.Time{}, ErrBootstrapInvalid
	}
	delete(m.bootstraps, hash)

	sessionID, err = randomToken(32)
	if err != nil {
		return "", "", time.Time{}, err
	}
	csrfToken, err = randomToken(32)
	if err != nil {
		return "", "", time.Time{}, err
	}
	expiresAt = m.now().Add(m.sessionTTL)
	m.sessions[hashToken(sessionID)] = sessionRecord{
		csrfHash:  hashToken(csrfToken),
		expiresAt: expiresAt,
	}
	return sessionID, csrfToken, expiresAt, nil
}

func (m *SessionManager) Validate(sessionID, csrfToken string, requireCSRF bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pruneLocked()
	record, ok := m.sessions[hashToken(sessionID)]
	if !ok || !record.expiresAt.After(m.now()) {
		return ErrSessionInvalid
	}
	if requireCSRF {
		provided := hashToken(csrfToken)
		if subtle.ConstantTimeCompare(record.csrfHash[:], provided[:]) != 1 {
			return ErrCSRFInvalid
		}
	}
	return nil
}

func (m *SessionManager) Revoke(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, hashToken(sessionID))
}

func (m *SessionManager) pruneLocked() {
	now := m.now()
	for token, record := range m.bootstraps {
		if !record.expiresAt.After(now) {
			delete(m.bootstraps, token)
		}
	}
	for token, record := range m.sessions {
		if !record.expiresAt.After(now) {
			delete(m.sessions, token)
		}
	}
}

func hashToken(token string) [sha256.Size]byte {
	return sha256.Sum256([]byte(token))
}

func randomToken(size int) (string, error) {
	buffer := make([]byte, size)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buffer), nil
}
