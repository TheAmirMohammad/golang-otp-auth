// internal/http/router.go (JWT middleware and protected routes)
package httpapi

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/golang-jwt/jwt/v5"

	"github.com/TheAmirMohammad/otp-service/internal/http/handlers"
)

func New(app *fiber.App, ah *handlers.AuthHandler, uh *handlers.UserHandler) {
	api := app.Group("/api/v1")

	//Auth endpoints
	api.Post("/auth/request-otp", ah.RequestOTP)
	api.Post("/auth/verify-otp", ah.VerifyOTP)
	
	protected := api.Group("", func(c *fiber.Ctx) error {
		h := c.Get("Authorization")
		if !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error":"missing bearer"})
		}
		tok := strings.TrimSpace(h[7:])
		t, err := jwt.Parse(tok, func(t *jwt.Token) (any, error) {
			return []byte(ah.JWTSecret), nil
		})
		if err != nil || !t.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error":"invalid token"})
		}
		return c.Next()
	})

	//User endpoints
	protected.Get("/users/:id", uh.GetUser)
	protected.Get("/users", uh.ListUsers)
	
	app.Get("/swagger/*", swagger.HandlerDefault)
}
