package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Cache CacheConf
	HTTP  HTTPConf
	App   AppConf
}

type AppConf struct {
	ServiceName string `env:"SERVICE_NAME" env-default:"reziser"`
	LogLevel    string `env:"LOG_LEVEL" env-default:"info"`
}
type CacheConf struct {
	MaxItems int    `env:"CACHE_ITEMS" env-default:"10"`
	Path     string `env:"CACHE_PATH" env-default:"./cache"`
}

type HTTPConf struct {
	Host string `env:"HOST" env-default:"0.0.0.0"`
	Port int    `env:"PORT" env-default:"8081"`
}

func Load() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	return &cfg, nil
}
