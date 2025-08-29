package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/TheAmirMohammad/otp-service/internal/config"
	httpapi "github.com/TheAmirMohammad/otp-service/internal/http"
	"github.com/TheAmirMohammad/otp-service/internal/http/handlers"
	"github.com/TheAmirMohammad/otp-service/internal/otp"
)

func main() {
	cfg := config.Load()
	app := fiber.New()
	otpMgr := otp.NewManager(2 * time.Minute)
	limiter := otp.NewLimiter(3, 10*time.Minute)

	ah := &handlers.AuthHandler{OTP: otpMgr, Limiter: limiter}

	app.Get("/health", func(c *fiber.Ctx) error { return c.SendString("OK") })
	httpapi.New(app, ah) //router

	log.Printf("listening on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil { log.Fatal(err) }
}
