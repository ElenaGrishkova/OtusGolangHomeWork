package main

import (
	"context"
	"os"

	//nolint:depguard
	"github.com/heetch/confita"
	//nolint:depguard
	"github.com/heetch/confita/backend"
	//nolint:depguard
	"github.com/heetch/confita/backend/env"
	//nolint:depguard
	"github.com/heetch/confita/backend/file"
	//nolint:depguard
	"github.com/heetch/confita/backend/flags"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Logger   LoggerConfig
}

type DatabaseConfig struct {
	Driver  string `config:"database-driver"`
	Dsn     string `config:"database-dsn"`
	Storage string `config:"database-storagetype"`
}

type ServerConfig struct {
	Host string `config:"server-host"`
	Port string `config:"server-port"`
}

type LoggerConfig struct {
	Level string `config:"logger-level"`
}

func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
		Database: DatabaseConfig{
			Storage: "in-memory",
		},
	}
}

func LoadConfig(ctx context.Context, cfg *Config, path string) error {
	backends := make([]backend.Backend, 0)
	if _, err := os.Stat(path); err == nil {
		backends = append(backends, file.NewBackend(path))
	}
	backends = append(backends, env.NewBackend())
	backends = append(backends, flags.NewBackend())

	loader := confita.NewLoader(backends...)
	if err := loader.Load(ctx, cfg); err != nil {
		return err
	}

	return nil
}
