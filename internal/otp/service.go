package otp

import "context"

// OTP service interface (both memory & redis implement)
type Service interface {
	Generate(ctx context.Context, phone string) (string, error)
	Validate(ctx context.Context, phone, code string) (bool, error)
}

// Rate limiter interface (both memory & redis implement)
type Limiter interface {
	Allow(ctx context.Context, phone string) (bool, error)
}
