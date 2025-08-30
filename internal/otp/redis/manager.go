package otp

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisManager struct {
	RDB *redis.Client
	TTL time.Duration
}

func NewRedisManager(rdb *redis.Client, ttl time.Duration) *RedisManager {
	return &RedisManager{RDB: rdb, TTL: ttl}
}

func (m *RedisManager) Generate(phone string) (string, error) {
	code, err := genCode()
	if err != nil { return "", err }
	key := m.key(phone)
	ctx := context.Background()

	// SET key code EX TTL NX (do not overwrite an unexpired OTP)
	ok, err := m.RDB.SetNX(ctx, key, code, m.TTL).Result()
	if err != nil { return "", err }
	if !ok {
		// Overwrite anyway: security-wise you may want to rotate; choose policy
		if err := m.RDB.Set(ctx, key, code, m.TTL).Err(); err != nil { return "", err }
	}
	log.Printf("[OTP] phone=%s code=%s (expires in %s)", phone, code, m.TTL)
	return code, nil
}

func (m *RedisManager) Validate(phone, code string) bool {
	ctx := context.Background()
	key := m.key(phone)
	val, err := m.RDB.Get(ctx, key).Result()
	if err != nil { return false } // includes redis.Nil
	if val != code { return false }
	// one-time use
	_ = m.RDB.Del(ctx, key).Err()
	return true
}

func (m *RedisManager) key(phone string) string { return "otp:code:" + phone }

func genCode() (string, error) {
	var b [3]byte
	if _, err := rand.Read(b[:]); err != nil { return "", err }
	n := (int(b[0])<<16 | int(b[1])<<8 | int(b[2])) % 1000000
	return fmt.Sprintf("%06d", n), nil
}
