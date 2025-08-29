package handlers

import (
	"context"
	"net/http"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/TheAmirMohammad/otp-service/internal/domain/user"
	jwtutil "github.com/TheAmirMohammad/otp-service/internal/jwt"
	"github.com/TheAmirMohammad/otp-service/internal/otp"
)

type AuthHandler struct {
	OTP      *otp.Manager
	Limiter  *otp.Limiter
	JWTSecret string
	TokenTTL  time.Duration
	Users     user.Repository
}

// DTOs (exported for Swagger)

type VerifyOTPReq struct {
	Phone string `json:"phone"`
	OTP   string `json:"otp"`
}
type AuthResp struct {
	Token string    `json:"token"`
	User  user.User `json:"user"`
}

type RequestOTPReq struct { Phone string `json:"phone"` }

var phoneRx = regexp.MustCompile(`^[0-9+\-() ]{5,20}$`) // For Iran numbers it should be "^09\d{9}$"

// RequestOTP godoc
// @Summary      Request OTP
// @Description  Generates an OTP (printed in server logs). Rate limit: 3 per 10 minutes per phone. Expires in 2 minutes.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload body RequestOTPReq true "Phone payload"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      429 {object} map[string]string
// @Router       /auth/request-otp [post]
func (h *AuthHandler) RequestOTP(c *fiber.Ctx) error {
	var req RequestOTPReq
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

// VerifyOTP godoc
// @Summary      Verify OTP (login/register)
// @Description  Validates OTP; creates user if not exists; returns JWT.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload body VerifyOTPReq true "Verify payload"
// @Success      200 {object} AuthResp
// @Failure      400 {object} map[string]string
// @Router       /auth/verify-otp [post]
func (h *AuthHandler) VerifyOTP(c *fiber.Ctx) error {
	var req VerifyOTPReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error":"invalid body"})
	}
	if !phoneRx.MatchString(req.Phone) || len(req.OTP) != 6 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error":"invalid inputs"})
	}
	if !h.OTP.Validate(req.Phone, req.OTP) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error":"invalid or expired otp"})
	}
	ctx := context.Background()
	u, _ := h.Users.GetByPhone(ctx, req.Phone)
	if u == nil {
		u = &user.User{ID: uuid.NewString(), Phone: req.Phone, RegisteredAt: time.Now().UTC()}
		if err := h.Users.Create(ctx, u); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error":"create user failed"})
		}
	}
	tok, err := jwtutil.Generate(h.JWTSecret, u.ID, h.TokenTTL)
	if err != nil { return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error":"token error"}) }
	return c.JSON(AuthResp{Token: tok, User: *u})
}
