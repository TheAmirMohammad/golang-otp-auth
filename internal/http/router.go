package httpapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	"github.com/TheAmirMohammad/otp-service/internal/http/handlers"
)

func New(app *fiber.App, ah *handlers.AuthHandler) {
	api := app.Group("/api/v1")
	api.Post("/auth/request-otp", ah.RequestOTP)
	api.Post("/auth/verify-otp", ah.VerifyOTP)
	
	app.Get("/swagger/*", swagger.HandlerDefault)
}
