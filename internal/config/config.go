package config

import "os"

type Config struct {
	Port      string
}

func Load() Config {
	return Config{
		Port:      env("PORT", "8080"),
	}
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" { return v }
	return d
}
