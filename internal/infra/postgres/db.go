package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Non-fatal connect: returns error instead of exiting
func Connect(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse db dsn: %w", err)
	}
	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}
	return db, nil
}

// Non-fatal migrate: returns error instead of exiting
func Migrate(ctx context.Context, db *pgxpool.Pool) error {
	const q = `
CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  phone TEXT NOT NULL UNIQUE,
  registered_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);`
	if _, err := db.Exec(ctx, q); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}
