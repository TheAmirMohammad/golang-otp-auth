package postgres

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func MustConnect(ctx context.Context, dsn string) *pgxpool.Pool {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("parse db dsn: %v", err)
	}
	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	return db
}

func MustMigrate(ctx context.Context, db *pgxpool.Pool) {
	const q = `
	CREATE TABLE IF NOT EXISTS users (
  	id TEXT PRIMARY KEY,
 	phone TEXT NOT NULL UNIQUE,
	registered_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	`
	if _, err := db.Exec(ctx, q); err != nil {
		log.Fatalf("migrate: %v", err)
	}
}
