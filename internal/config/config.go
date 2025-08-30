package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	JWTSecret string

	// Toggles (single .env friendly)
	UseDB    bool
	UseRedis bool

	// Effective URLs used by the app (empty => in-memory)
	DatabaseURL string
	RedisURL    string

	// Base pieces (for building URLs if URLs not provided)
	PGUser     string
	PGPassword string
	PGDB       string
	PGHost     string
	PGPort     int

	RedisHost string
	RedisPort int
	RedisDB   int
}

func Load() Config {
	// Try to load .env; warn (don’t crash) if not found
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: .env not found or unreadable (%v) – continuing with process env", err)
	}

	cfg := Config{
		Port:      env("PORT", "8080"),
		JWTSecret: env("JWT_SECRET", "golangotpauthentication"),

		UseDB:    envBool("USE_DB", true),
		UseRedis: envBool("USE_REDIS", true),

		DatabaseURL: strings.TrimSpace(os.Getenv("DATABASE_URL")),
		RedisURL:    strings.TrimSpace(os.Getenv("REDIS_URL")),

		PGUser:     env("POSTGRES_USER", "otp"),
		PGPassword: env("POSTGRES_PASSWORD", "otp"),
		PGDB:       env("POSTGRES_DB", "otp"),
		PGHost:     env("POSTGRES_DNS", "db"),
		PGPort:     envInt("POSTGRES_PORT", 5432),

		RedisHost: env("REDIS_DNS", "redis"),
		RedisPort: envInt("REDIS_PORT", 6379),
		RedisDB:   envInt("REDIS_DB", 0),
	}

	// Toggles force in-memory by blanking URLs
	if !cfg.UseDB {
		cfg.DatabaseURL = ""
	}
	if !cfg.UseRedis {
		cfg.RedisURL = ""
	}

	// Build URLs from pieces if toggled on but URL empty
	if cfg.UseDB && cfg.DatabaseURL == "" {
		cfg.DatabaseURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable",
			urlEscape(cfg.PGUser), urlEscape(cfg.PGPassword),
			cfg.PGHost, cfg.PGPort, urlEscape(cfg.PGDB),
		)
	}
	if cfg.UseRedis && cfg.RedisURL == "" {
		cfg.RedisURL = fmt.Sprintf("redis://%s:%d/%d", cfg.RedisHost, cfg.RedisPort, cfg.RedisDB)
	}

	return cfg
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func envInt(k string, d int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return d
}

func envBool(k string, d bool) bool {
	if v := strings.TrimSpace(os.Getenv(k)); v != "" {
		switch strings.ToLower(v) {
		case "1", "true", "t", "yes", "y", "on":
			return true
		case "0", "false", "f", "no", "n", "off":
			return false
		}
	}
	return d
}

// Minimal percent-escape for URL segments
func urlEscape(s string) string {
	r := strings.NewReplacer(" ", "%20", "#", "%23", "@", "%40", ":", "%3A", "/", "%2F", "?", "%3F", "&", "%26", "=", "%3D")
	return r.Replace(s)
}
