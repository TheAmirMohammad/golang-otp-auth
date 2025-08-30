package memoryotp

import (
	"context"
	"sync"
	"time"

	"github.com/TheAmirMohammad/otp-service/internal/otp"
)

type limiter struct {
	mu      sync.Mutex
	window  time.Duration
	limit   int
	records map[string][]time.Time
}

func NewLimiter(limit int, window time.Duration) otp.Limiter {
	return &limiter{
		window:  window,
		limit:   limit,
		records: make(map[string][]time.Time),
	}
}

func (l *limiter) Allow(_ context.Context, phone string) (bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	cut := now.Add(-l.window)
	arr := l.records[phone][:0]
	for _, t := range l.records[phone] {
		if t.After(cut) {
			arr = append(arr, t)
		}
	}
	if len(arr) >= l.limit {
		l.records[phone] = arr
		return false, nil
	}
	l.records[phone] = append(arr, now)
	return true, nil
}
