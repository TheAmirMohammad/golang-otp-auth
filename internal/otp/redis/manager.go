package redisotp

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	
	"github.com/TheAmirMohammad/otp-service/internal/otp"
)

type manager struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewManager(rdb *redis.Client, ttl time.Duration) otp.Service {
	return &manager{rdb: rdb, ttl: ttl}
}

func (m *manager) Generate(ctx context.Context, phone string) (string, error) {
	code, err := genCode()
	if err != nil {
		return "", err
	}
	key := fmt.Sprintf("otp:%s", phone)
	if err := m.rdb.Set(ctx, key, code, m.ttl).Err(); err != nil {
		return "", err
	}
	log.Printf("[OTP] phone=%s code=%s (expires in %s)", phone, code, m.ttl)
	return code, nil
}

func (m *manager) Validate(ctx context.Context, phone, code string) (bool, error) {
	key := fmt.Sprintf("otp:%s", phone)
	val, err := m.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if val != code {
		return false, nil
	}
	if err := m.rdb.Del(ctx, key).Err(); err != nil {
		return false, err
	}
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
