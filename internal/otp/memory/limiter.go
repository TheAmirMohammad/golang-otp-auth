package otp

import (
	"sync"
	"time"
)

type Limiter struct {
	mu      sync.Mutex
	window  time.Duration
	limit   int
	records map[string][]time.Time
}

func NewLimiter(limit int, window time.Duration) *Limiter {
	return &Limiter{limit: limit, window: window, records: map[string][]time.Time{}}
}

func (l *Limiter) Allow(key string) bool {
	l.mu.Lock(); defer l.mu.Unlock()
	now := time.Now()
	cut := now.Add(-l.window)
	arr := l.records[key][:0]
	for _, t := range l.records[key] {
		if t.After(cut) { arr = append(arr, t) }
	}
	if len(arr) >= l.limit {
		l.records[key] = arr
		return false
	}
	l.records[key] = append(arr, now)
	return true
}
