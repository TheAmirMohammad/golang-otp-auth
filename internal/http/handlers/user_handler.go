package handlers

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/TheAmirMohammad/otp-service/internal/domain/user"
)

type UserHandler struct { Users user.Repository }

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	u, _ := h.Users.GetByID(context.Background(), id)
	if u == nil { return c.Status(http.StatusNotFound).JSON(fiber.Map{"error":"not found"}) }
	return c.JSON(u)
}