package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env         string
	Port        string
	DatabaseURL string
	RedisAddr   string
	JWTSecret   string
}

func Load() *Config {
	_ = godotenv.Load()
	return &Config{
		Env:         env("ENV", "dev"),
		Port:        env("PORT", "8081"),
		DatabaseURL: must("DATABASE_URL"),
		RedisAddr:   env("REDIS_ADDR", "localhost:6379"),
		JWTSecret:   env("JWT_SECRET", "change_me_dev"),
	}
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" { return v }
	return d
}
func must(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing required env %s", k)
	}
	return v
}

var DefaultTimeout = 10 * time.Second
