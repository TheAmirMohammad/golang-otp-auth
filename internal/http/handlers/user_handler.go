package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/TheAmirMohammad/otp-service/internal/domain/user"
)

type UserHandler struct { Users user.Repository }

type listResp struct {
	Items []user.User `json:"items"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	u, _ := h.Users.GetByID(context.Background(), id)
	if u == nil { return c.Status(http.StatusNotFound).JSON(fiber.Map{"error":"not found"}) }
	return c.JSON(u)
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	size, _ := strconv.Atoi(c.Query("size", "20"))
	if page < 1 { page = 1 }
	if size < 1 || size > 100 { size = 20 }
	search := c.Query("search", "")
	items, total, err := h.Users.List(context.Background(), user.ListFilter{
		Search: search, Limit: size, Offset: (page-1)*size,
	})
	if err != nil { return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error":"list error"}) }
	return c.JSON(listResp{Items: items, Total: total, Page: page, Size: size})
}