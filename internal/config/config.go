package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	JWTSecret string

	// Toggles
	UseDB    bool
	UseRedis bool

	// Effective URLs (empty => in-memory)
	DatabaseURL string
	RedisURL    string

	// Base pieces to build URLs
	PGUser     string
	PGPassword string
	PGDB       string
	PGHost     string
	PGPort     int

	RedisHost string
	RedisPort int
	RedisDB   int

	// ⚙️ Tunables
	OTPTTL         time.Duration // default 2m
	RateLimitMax   int           // default 3
	RateLimitWindow time.Duration // default 10m
	TokenTTL       time.Duration // default 24h
}

func Load() Config {
	// Load .env if present (warn if missing)
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

		// Tunables (durations accept Go format: 30s, 2m, 1h)
		OTPTTL:          envDuration("OTP_TTL", 2*time.Minute),
		RateLimitMax:    envInt("RATE_LIMIT_MAX", 3),
		RateLimitWindow: envDuration("RATE_LIMIT_WINDOW", 10*time.Minute),
		TokenTTL:        envDuration("TOKEN_TTL", 24*time.Hour),
	}

	// Toggles force in-memory by blanking URLs
	if !cfg.UseDB {
		cfg.DatabaseURL = ""
	}
	if !cfg.UseRedis {
		cfg.RedisURL = ""
	}

	// Build URLs from pieces if needed
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
func envDuration(k string, d time.Duration) time.Duration {
	if v, ok := os.LookupEnv(k); ok && strings.TrimSpace(v) != "" {
		if dur, err := time.ParseDuration(strings.TrimSpace(v)); err == nil {
			return dur
		}
		log.Printf("warning: invalid duration for %s=%q (use e.g. 30s, 2m, 1h); using default %s", k, v, d)
	}
	return d
}

// Minimal percent-escape for URL segments
func urlEscape(s string) string {
	r := strings.NewReplacer(" ", "%20", "#", "%23", "@", "%40", ":", "%3A", "/", "%2F", "?", "%3F", "&", "%26", "=", "%3D")
	return r.Replace(s)
}
