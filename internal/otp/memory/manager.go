package memoryotp

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/TheAmirMohammad/otp-service/internal/otp"
)

type manager struct {
	mu  sync.RWMutex
	ttl time.Duration
	m   map[string]record
}

type record struct {
	Code      string
	ExpiresAt time.Time
}

func NewManager(ttl time.Duration) otp.Service {
	return &manager{ttl: ttl, m: make(map[string]record)}
}

func (m *manager) Generate(_ context.Context, phone string) (string, error) {
	code, err := genCode()
	if err != nil {
		return "", err
	}
	m.mu.Lock()
	m.m[phone] = record{Code: code, ExpiresAt: time.Now().Add(m.ttl)}
	m.mu.Unlock()
	log.Printf("[OTP] phone=%s code=%s (expires in %s)", phone, code, m.ttl)
	return code, nil
}

func (m *manager) Validate(_ context.Context, phone, code string) (bool, error) {
	m.mu.RLock()
	rec, ok := m.m[phone]
	m.mu.RUnlock()
	if !ok || time.Now().After(rec.ExpiresAt) || rec.Code != code {
		return false, nil
	}
	m.mu.Lock()
	delete(m.m, phone) // one-time use
	m.mu.Unlock()
	return true, nil
}

func genCode() (string, error) {
	var b [3]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	n := (int(b[0])<<16 | int(b[1])<<8 | int(b[2])) % 1000000
	return fmt.Sprintf("%06d", n), nil
}
