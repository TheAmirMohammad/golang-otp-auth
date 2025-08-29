package memory

import (
	"sync"

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
