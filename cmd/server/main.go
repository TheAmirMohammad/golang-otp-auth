package main

import (
	"context"
	"log"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	_ "github.com/TheAmirMohammad/otp-service/docs" // swagger docs

	"github.com/TheAmirMohammad/otp-service/internal/config"
	"github.com/TheAmirMohammad/otp-service/internal/domain/user"
	httpapi "github.com/TheAmirMohammad/otp-service/internal/http"
	"github.com/TheAmirMohammad/otp-service/internal/http/handlers"
	"github.com/TheAmirMohammad/otp-service/internal/infra/memory"
	"github.com/TheAmirMohammad/otp-service/internal/infra/postgres"
	"github.com/TheAmirMohammad/otp-service/internal/otp"
	mem "github.com/TheAmirMohammad/otp-service/internal/otp/memory"
	red "github.com/TheAmirMohammad/otp-service/internal/otp/redis"
)

// @title           OTP Service API
// @version         1.0
// @description     OTP-based login/registration with user management.
// @BasePath        /api/v1
// @securityDefinitions.apikey Bearer
// @in              header
// @name            Authorization
func main() {
	ctx := context.Background()
	cfg := config.Load()

	usersRepo := buildUserRepo(ctx, cfg)
	otpSvc, limiter := buildOTPStack(ctx, cfg)

	ah := &handlers.AuthHandler{
		Users:     usersRepo,
		OTP:       otpSvc,
		Limiter:   limiter,
		JWTSecret: cfg.JWTSecret,
		TokenTTL:  cfg.TokenTTL,
	}
	uh := &handlers.UserHandler{Users: usersRepo}

	app := fiber.New()
	app.Get("/health", func(c *fiber.Ctx) error { return c.SendString("OK") })
	httpapi.New(app, ah, uh)

	log.Printf("config: otp_ttl=%v rate=%d/%v token_ttl=%v", cfg.OTPTTL, cfg.RateLimitMax, cfg.RateLimitWindow, cfg.TokenTTL)
	log.Printf("listening on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}

// buildUserRepo wires Postgres if available, otherwise falls back to memory.
func buildUserRepo(ctx context.Context, cfg config.Config) user.Repository {
	if strings.TrimSpace(cfg.DatabaseURL) == "" {
		log.Println("users repo: in-memory (DATABASE_URL empty or USE_DB=false)")
		return memory.NewUserRepo()
	}
	db, err := postgres.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Printf("warning: postgres unavailable (%v) – using in-memory users repo", err)
		return memory.NewUserRepo()
	}
	if err := postgres.Migrate(ctx, db); err != nil {
		log.Printf("warning: migration failed (%v) – using in-memory users repo", err)
		return memory.NewUserRepo()
	}
	log.Println("users repo: postgres")
	return postgres.NewUserRepo(db)
}

// buildOTPStack wires Redis-backed OTP & rate if available, otherwise falls back to memory.
func buildOTPStack(ctx context.Context, cfg config.Config) (otp.Service, otp.Limiter) {
	if strings.TrimSpace(cfg.RedisURL) == "" {
		log.Println("otp/rate: in-memory (REDIS_URL empty or USE_REDIS=false)")
		return mem.NewManager(cfg.OTPTTL), mem.NewLimiter(cfg.RateLimitMax, cfg.RateLimitWindow)
	}
	rdb := redis.NewClient(mustParseRedisURL(cfg.RedisURL))
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("warning: redis unavailable (%v) – using in-memory OTP & rate", err)
		return mem.NewManager(cfg.OTPTTL), mem.NewLimiter(cfg.RateLimitMax, cfg.RateLimitWindow)
	}
	log.Println("otp/rate: redis")
	return red.NewManager(rdb, cfg.OTPTTL), red.NewLimiter(rdb, cfg.RateLimitMax, cfg.RateLimitWindow)
}

// mustParseRedisURL accepts either "host:port" or "redis://[:pass@]host:port[/db]"
func mustParseRedisURL(raw string) *redis.Options {
	if !strings.Contains(raw, "://") {
		return &redis.Options{Addr: raw}
	}
	u, err := url.Parse(raw)
	if err != nil {
		log.Fatalf("parse redis url: %v", err)
	}
	if u.Scheme != "redis" && u.Scheme != "rediss" {
		log.Fatalf("unsupported redis scheme: %s", u.Scheme)
	}
	opts, err := redis.ParseURL(raw)
	if err != nil {
		log.Fatalf("parse redis url (ParseURL): %v", err)
	}
	return opts
}
