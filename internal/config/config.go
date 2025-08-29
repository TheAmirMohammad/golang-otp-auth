package config

import "os"

type Config struct {
	Port      string
	JWTSecret string
}

func Load() Config {
	return Config{
		Port:      env("PORT", "8080"),
		JWTSecret: env("JWT_SECRET", "golnagotpauthentication"),
	}
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" { return v }
	return d
}