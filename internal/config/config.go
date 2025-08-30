package config

import "os"

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseURL string // if empty => memory users repo
	RedisURL    string // if empty => in-memory OTP & rate
}

func Load() Config {
	return Config{
		Port:        env("PORT", "8080"),
		JWTSecret:   env("JWT_SECRET", "golangotpauthentication"),
		DatabaseURL: env("DATABASE_URL", ""),
		RedisURL:    env("REDIS_URL", ""),
	}
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" { return v }
	return d
}