package config

import (
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string `env:"APP_ENV"`
	AppPort string `env:"APP_PORT" envDefault:"8080"`

	DBHost      string `env:"DB_HOST"`
	DBPort      string `env:"DB_PORT"`
	DBUser      string `env:"DB_USER"`
	DBPassword  string `env:"DB_PASSWORD"`
	DBName      string `env:"DB_NAME"`
	DatabaseURL string `env:"DATABASE_URL"`

	RedisAddr     string `env:"REDIS_ADDR"`
	RedisPassword string `env:"REDIS_PASSWORD"`

	FXAPIURL string `env:"FX_API_URL"`
	FXAPIKey string `env:"FX_API_KEY"`

	CORSAllowedOrigins string `env:"CORS_ALLOWED_ORIGINS"`

	RateLimitRPS   int `env:"RATE_LIMIT_RPS"`
	RateLimitBurst int `env:"RATE_LIMIT_BURST"`
}

func Load() *Config {
	_ = godotenv.Load()
	cfg := &Config{}
	_ = env.Parse(cfg)
	return cfg
}

func (c *Config) GetCORSOrigins() []string {
	if c.CORSAllowedOrigins == "" {
		return []string{}
	}
	parts := strings.Split(c.CORSAllowedOrigins, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
