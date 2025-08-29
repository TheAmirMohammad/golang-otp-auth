package handlers

import (
	"regexp"

	"github.com/TheAmirMohammad/otp-service/internal/otp"
)

type AuthHandler struct {
	OTP      *otp.Manager
	Limiter  *otp.Limiter
}

type requestOTPReq struct { Phone string `json:"phone"` }

var phoneRx = regexp.MustCompile(`^[0-9+\-() ]{5,20}$`) // For Iran numbers it should be "^09\d{9}$"
