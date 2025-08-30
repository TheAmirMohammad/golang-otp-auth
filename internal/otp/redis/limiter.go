package redisotp

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	
	"github.com/TheAmirMohammad/otp-service/internal/otp"
)

type limiter struct {
	rdb    *redis.Client
	limit  int
	window time.Duration
}

func NewLimiter(rdb *redis.Client, limit int, window time.Duration) otp.Limiter {
	return &limiter{rdb: rdb, limit: limit, window: window}
}

func (l *limiter) Allow(ctx context.Context, phone string) (bool, error) {
	key := fmt.Sprintf("rl:otp:%s", phone)
	n, err := l.rdb.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if n == 1 {
		_ = l.rdb.Expire(ctx, key, l.window).Err()
	}
	return n <= int64(l.limit), nil
}
