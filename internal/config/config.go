package config

import (
	"log/slog"

	slerr "github.com/RakhimovAns/logger/pkg/err"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	TelegramToken string `env:"TELEGRAM_TOKEN"`

	PostgresHost          string `env:"POSTGRES_HOST"`
	PostgresPort          string `env:"POSTGRES_PORT"`
	PostgresUser          string `env:"POSTGRES_USER"`
	PostgresDB            string `env:"POSTGRES_DB"`
	PostgresPassword      string `env:"POSTGRES_PASSWORD"`
	PostgresMigrationPass string `env:"POSTGRES_MIGRATION_PASS"`

	RedisHost     string `env:"REDIS_HOST"`
	RedisPort     string `env:"REDIS_PORT"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       int    `env:"REDIS_DB"`
}

func NewConfig() *Config {
	cfg := Config{}
	err := cleanenv.ReadConfig("./configs/local.env", &cfg)
	if err != nil {
		slog.Error(slerr.WithSource(err).Error())
		return nil
	}
	slog.Info("Config loaded")
	return &cfg
}

var cfg = NewConfig()
