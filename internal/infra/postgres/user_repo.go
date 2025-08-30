package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/TheAmirMohammad/otp-service/internal/domain/user"
)

type UserRepo struct{ db *pgxpool.Pool }

func NewUserRepo(db *pgxpool.Pool) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) Create(ctx context.Context, u *user.User) error {
	_, err := r.db.Exec(ctx, `INSERT INTO users (id, phone, registered_at) VALUES ($1,$2,$3)`,
		u.ID, u.Phone, u.RegisteredAt)
	return err
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*user.User, error) {
	row := r.db.QueryRow(ctx, `SELECT id, phone, registered_at FROM users WHERE id=$1`, id)
	var u user.User
	if err := row.Scan(&u.ID, &u.Phone, &u.RegisteredAt); err != nil {
		return nil, nil
	}
	return &u, nil
}

func (r *UserRepo) GetByPhone(ctx context.Context, phone string) (*user.User, error) {
	row := r.db.QueryRow(ctx, `SELECT id, phone, registered_at FROM users WHERE phone=$1`, phone)
	var u user.User
	if err := row.Scan(&u.ID, &u.Phone, &u.RegisteredAt); err != nil {
		return nil, nil
	}
	return &u, nil
}

func (r *UserRepo) List(ctx context.Context, f user.ListFilter) ([]user.User, int, error) {
	var (
		q    = `SELECT id, phone, registered_at FROM users`
		args []any
	)
	if s := strings.TrimSpace(f.Search); s != "" {
		q += ` WHERE phone ILIKE $1`
		args = append(args, "%"+s+"%")
	}
	// LIMIT/OFFSET placeholders based on current args length
	args = append(args, f.Limit, f.Offset)
	q += fmt.Sprintf(` ORDER BY registered_at DESC LIMIT $%d OFFSET $%d`, len(args)-1, len(args))

	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.ID, &u.Phone, &u.RegisteredAt); err != nil {
			return nil, 0, err
		}
		out = append(out, u)
	}

	// total
	var total int
	if strings.TrimSpace(f.Search) != "" {
		if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE phone ILIKE $1`, "%"+f.Search+"%").Scan(&total); err != nil {
			return nil, 0, err
		}
	} else {
		if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total); err != nil {
			return nil, 0, err
		}
	}
	return out, total, nil
}
