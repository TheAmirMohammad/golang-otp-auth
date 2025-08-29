package otp

import (
	"crypto/rand"
	"fmt"
	"log"
	"sync"
	"time"
)


type record struct {
	Code      string
	ExpiresAt time.Time
}

type Manager struct {
	mu    sync.RWMutex
	store map[string]record
	ttl   time.Duration
}

func NewManager(ttl time.Duration) *Manager {
	m := &Manager{store: map[string]record{}, ttl: ttl}
	// best-effort cleanup
	go func() {
		t := time.NewTicker(time.Minute)
		for range t.C { m.cleanup() }
	}()
	return m
}

func (m *Manager) cleanup() {
	now := time.Now()
	m.mu.Lock(); defer m.mu.Unlock()
	for k, v := range m.store {
		if now.After(v.ExpiresAt) { delete(m.store, k) }
	}
}

func genCode() (string, error) {
	var b [3]byte
	if _, err := rand.Read(b[:]); err != nil { return "", err }
	n := (int(b[0])<<16 | int(b[1])<<8 | int(b[2])) % 1000000
	return fmt.Sprintf("%06d", n), nil
}

func (m *Manager) Generate(phone string) (string, error) {
	code, err := genCode()
	if err != nil { return "", err }
	m.mu.Lock()
	m.store[phone] = record{Code: code, ExpiresAt: time.Now().Add(m.ttl)}
	log.Printf("[OTP] %s -> %s (expires in %s)", phone, code, m.ttl)
	return code, nil
}