package otp

import (
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