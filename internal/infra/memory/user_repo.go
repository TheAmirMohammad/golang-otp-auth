package memory

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/TheAmirMohammad/otp-service/internal/domain/user"
)

type UserRepo struct {
	mu      sync.RWMutex
	byID    map[string]user.User
	byPhone map[string]string
}

func NewUserRepo() *UserRepo {
	return &UserRepo{byID: map[string]user.User{}, byPhone: map[string]string{}}
}

func (r *UserRepo) Create(ctx context.Context, u *user.User) error {
	r.mu.Lock(); defer r.mu.Unlock()
	if u.ID == "" { u.ID = uuid.NewString() }
	if u.RegisteredAt.IsZero() { u.RegisteredAt = time.Now().UTC() }
	r.byID[u.ID] = *u
	r.byPhone[u.Phone] = u.ID
	return nil
}
