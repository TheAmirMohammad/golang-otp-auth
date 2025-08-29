package httpapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func New(app *fiber.App) {
	// api := app.Group("/api/v1")
	app.Get("/swagger/*", swagger.HandlerDefault)
}
