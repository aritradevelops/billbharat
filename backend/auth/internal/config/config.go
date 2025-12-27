package config

import (
	"fmt"
	"time"

	// autoload the environment variables

	"github.com/caarlos0/env/v10"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Http       Http       `envPrefix:"HTTP_"`
	Database   Database   `envPrefix:"DATABASE_"`
	Service    Service    `envPrefix:"SERVICE_"`
	Grpc       Grpc       `envPrefix:"GRPC_"`
	Deployment Deployment `envPrefix:"DEPLOYMENT_"`
}

type Http struct {
	Host string `env:"HOST,required"`
	Port int    `env:"PORT,required"`
}

type Database struct {
	Uri     string        `env:"URI,required"`
	Timeout time.Duration `env:"TIMEOUT,required"`
}

type Service struct {
	Name string `env:"NAME,required"`
}

type Grpc struct {
	Host string `env:"HOST,required"`
	Port int    `env:"PORT,required"`
}

type Deployment struct {
	Env string `env:"ENV,required"`
}

func Load() (Config, error) {
	var config Config
	err := env.Parse(&config)
	if err != nil {
		return config, fmt.Errorf("failed to load config: %v", err)
	}
	return config, nil
}
