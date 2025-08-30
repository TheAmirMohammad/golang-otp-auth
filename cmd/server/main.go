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

	// --- Users repo: Postgres if DATABASE_URL set, otherwise in-memory
	var usersRepo user.Repository
	if strings.TrimSpace(cfg.DatabaseURL) != "" {
		db := postgres.MustConnect(ctx, cfg.DatabaseURL)
		postgres.MustMigrate(ctx, db)
		usersRepo = postgres.NewUserRepo(db)
		log.Println("users repo: postgres")
	} else {
		usersRepo = memory.NewUserRepo()
		log.Println("users repo: in-memory")
	}

	// --- OTP & Rate: Redis if REDIS_URL set, otherwise in-memory
	var otpSvc otp.Service
	var limiter otp.Limiter

	if cfg.RedisURL != "" {
		rdb := redis.NewClient(mustParseRedisURL(cfg.RedisURL))
		otpSvc = red.NewManager(rdb, 2*time.Minute)
		limiter = red.NewLimiter(rdb, 3, 10*time.Minute)
		log.Println("otp: redis")
	} else {
		otpSvc = mem.NewManager(2 * time.Minute)
		limiter = mem.NewLimiter(3, 10*time.Minute)
		log.Println("otp: in-memory")
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
	app.Get("/health", func(c *fiber.Ctx) error { return c.SendString("OK") })
	httpapi.New(app, ah, uh)

	log.Printf("listening on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}

// mustParseRedisURL accepts either "host:port" or "redis://[:pass@]host:port[/db]"
func mustParseRedisURL(raw string) *redis.Options {
	// plain host:port
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
