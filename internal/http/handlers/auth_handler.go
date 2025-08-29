package handlers

import "github.com/TheAmirMohammad/otp-service/internal/otp"

type AuthHandler struct {
	OTP      *otp.Manager
	Limiter  *otp.Limiter
}

type requestOTPReq struct { Phone string `json:"phone"` }
