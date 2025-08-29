package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/TheAmirMohammad/otp-service/internal/config"
	httpapi "github.com/TheAmirMohammad/otp-service/internal/http"
	"github.com/TheAmirMohammad/otp-service/internal/http/handlers"
	"github.com/TheAmirMohammad/otp-service/internal/infra/memory"
	"github.com/TheAmirMohammad/otp-service/internal/otp"
	_ "github.com/TheAmirMohammad/otp-service/docs" // swagger docs
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
	app := fiber.New()

	repo := memory.NewUserRepo()
	otpMgr := otp.NewManager(2 * time.Minute)
	limiter := otp.NewLimiter(3, 10*time.Minute)

	ah := &handlers.AuthHandler{
		OTP: otpMgr, Limiter: limiter, Users: repo,
		JWTSecret: cfg.JWTSecret, TokenTTL: 24 * time.Hour,
	}

	app.Get("/health", func(c *fiber.Ctx) error { return c.SendString("OK") })
	httpapi.New(app, ah) //router

	log.Printf("listening on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil { log.Fatal(err) }
}
