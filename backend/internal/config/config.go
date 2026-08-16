package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Config holds all runtime configuration injected via environment variables.
type Config struct {
	Port          int    `env:"PORT" envDefault:"8080"`
	GinMode       string `env:"GIN_MODE" envDefault:"release"`
	DBHost        string `env:"DB_HOST" envDefault:"127.0.0.1"`
	DBPort        string `env:"DB_PORT" envDefault:"3306"`
	DBUser        string `env:"DB_USER" envDefault:"gbkanban"`
	DBPassword    string `env:"DB_PASSWORD" envDefault:"gbkanban"`
	DBName        string `env:"DB_NAME" envDefault:"gbkanban"`
	JWTSecret     string `env:"JWT_SECRET" envDefault:"gbkanban-secret"`
	TokenTTLHours int    `env:"JWT_TTL_HOURS" envDefault:"168"`
	UploadDir     string `env:"UPLOAD_DIR" envDefault:"/app/uploads"`
	AllowedOrigin string `env:"ALLOWED_ORIGIN" envDefault:"*"`
}

// Load parses configuration from the process environment.
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}

// DSN returns the MySQL data source name used by GORM.
func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}
