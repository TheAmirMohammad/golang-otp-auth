package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/TheAmirMohammad/otp-service/internal/config"
	httpapi "github.com/TheAmirMohammad/otp-service/internal/http"
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

	app.Get("/health", func(c *fiber.Ctx) error { return c.SendString("OK") })
	httpapi.New(app) //router

	log.Printf("listening on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil { log.Fatal(err) }
}
