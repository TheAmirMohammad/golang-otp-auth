package otp

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLimiter struct {
	RDB    *redis.Client
	Limit  int
	Window time.Duration
}

func NewRedisLimiter(rdb *redis.Client, limit int, window time.Duration) *RedisLimiter {
	return &RedisLimiter{RDB: rdb, Limit: limit, Window: window}
}

// Fixed-window counter: INCR key, set TTL on first hit, allow if <= Limit.
func (l *RedisLimiter) Allow(phone string) bool {
	ctx := context.Background()
	key := l.key(phone, time.Now(), l.Window)

	pipe := l.RDB.TxPipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, l.Window)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false // fail-closed or choose to allow on redis error
	}
	return incr.Val() <= int64(l.Limit)
}

func (l *RedisLimiter) key(phone string, now time.Time, window time.Duration) string {
	// window bucket start (fixed window)
	bucket := now.Unix() / int64(window.Seconds())
	return fmt.Sprintf("otp:rl:%s:%d", phone, bucket)
}
