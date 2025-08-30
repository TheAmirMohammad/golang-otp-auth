package main

import (
	"context"
	"log"
	"net/url"
	"strings"
	"time"

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
	cfg := config.Load()
	ctx := context.Background()

	// --- Users repo: try Postgres; fall back to memory on ANY error
	var usersRepo user.Repository
	if strings.TrimSpace(cfg.DatabaseURL) != "" {
		if db, err := postgres.Connect(ctx, cfg.DatabaseURL); err != nil {
			log.Printf("warning: postgres unavailable (%v) – using in-memory users repo", err)
			usersRepo = memory.NewUserRepo()
		} else if err := postgres.Migrate(ctx, db); err != nil {
			log.Printf("warning: migration failed (%v) – using in-memory users repo", err)
			usersRepo = memory.NewUserRepo()
		} else {
			usersRepo = postgres.NewUserRepo(db)
			log.Println("users repo: postgres")
		}
	} else {
		usersRepo = memory.NewUserRepo()
		log.Println("users repo: in-memory (DATABASE_URL empty or USE_DB=false)")
	}

	// --- OTP & Rate: Redis if REDIS_URL set; otherwise in-memory
	var otpSvc otp.Service
	var limiter otp.Limiter

	if strings.TrimSpace(cfg.RedisURL) != "" {
		rdb := redis.NewClient(mustParseRedisURL(cfg.RedisURL))
		if err := rdb.Ping(ctx).Err(); err != nil {
			log.Printf("warning: redis unavailable (%v) – using in-memory OTP & rate", err)
			otpSvc = mem.NewManager(2 * time.Minute)
			limiter = mem.NewLimiter(3, 10*time.Minute)
		} else {
			otpSvc = red.NewManager(rdb, 2*time.Minute)
			limiter = red.NewLimiter(rdb, 3, 10*time.Minute)
			log.Println("otp/rate: redis")
		}
	} else {
		otpSvc = mem.NewManager(2 * time.Minute)
		limiter = mem.NewLimiter(3, 10*time.Minute)
		log.Println("otp/rate: in-memory (REDIS_URL empty or USE_REDIS=false)")
	}

	// Handlers
	ah := &handlers.AuthHandler{
		Users:     usersRepo,
		OTP:       otpSvc,
		Limiter:   limiter,
		JWTSecret: cfg.JWTSecret,
		TokenTTL:  24 * time.Hour,
	}
	uh := &handlers.UserHandler{Users: usersRepo}

	app := fiber.New()
	// Simple health check
	app.Get("/health", func(c *fiber.Ctx) error { return c.SendString("OK") })
	httpapi.New(app, ah, uh)

	log.Printf("listening on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
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
