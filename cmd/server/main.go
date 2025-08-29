package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/TheAmirMohammad/otp-service/internal/config"
)

func main() {
	cfg := config.Load()
	app := fiber.New()
	app.Get("/health", func(c *fiber.Ctx) error { return c.SendString("OK") })
	log.Printf("listening on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil { log.Fatal(err) }
}
