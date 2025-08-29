package user

import "context"

type ListFilter struct {
	Search string
	Limit  int
	Offset int
}

type Repository interface {
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByPhone(ctx context.Context, phone string) (*User, error)
	List(ctx context.Context, f ListFilter) (users []User, total int, err error)
}
