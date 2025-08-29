package memory

import (
	"context"
	"strings"
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

func (r *UserRepo) GetByID(ctx context.Context, id string) (*user.User, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	u, ok := r.byID[id]
	if !ok { return nil, nil }
	return &u, nil
}

func (r *UserRepo) GetByPhone(ctx context.Context, phone string) (*user.User, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	id, ok := r.byPhone[phone]
	if !ok { return nil, nil }
	u := r.byID[id]
	return &u, nil
}

func (r *UserRepo) List(ctx context.Context, f user.ListFilter) ([]user.User, int, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	var out []user.User
	needle := strings.ToLower(strings.TrimSpace(f.Search))
	for _, u := range r.byID {
		if needle == "" || strings.Contains(strings.ToLower(u.Phone), needle) {
			out = append(out, u)
		}
	}
	total := len(out)
	start := min(f.Offset, total)
	end := min(start + f.Limit, total)
	return out[start:end], total, nil
}