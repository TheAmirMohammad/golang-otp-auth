package handlers

import (
	"net/http"
	"regexp"
	"time"

	"github.com/TheAmirMohammad/otp-service/internal/domain/user"
	"github.com/TheAmirMohammad/otp-service/internal/otp"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	OTP      *otp.Manager
	Limiter  *otp.Limiter
	JWTSecret string
	TokenTTL  time.Duration
	Users     user.Repository
}

type verifyOTPReq struct {
	Phone string `json:"phone"`
	OTP   string `json:"otp"`
}
type authResp struct {
	Token string    `json:"token"`
	User  user.User `json:"user"`
}

type requestOTPReq struct { Phone string `json:"phone"` }

var phoneRx = regexp.MustCompile(`^[0-9+\-() ]{5,20}$`) // For Iran numbers it should be "^09\d{9}$"

func (h *AuthHandler) RequestOTP(c *fiber.Ctx) error {
	var req requestOTPReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error":"invalid body"})
	}
	if !phoneRx.MatchString(req.Phone) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error":"invalid phone"})
	}
	if !h.Limiter.Allow(req.Phone) {
		return c.Status(http.StatusTooManyRequests).JSON(fiber.Map{"error":"rate limit: 3 per 10 minutes"})
	}
	if _, err := h.OTP.Generate(req.Phone); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error":"otp error"})
	}
	return c.JSON(fiber.Map{"message":"otp generated (check server logs)"})
}