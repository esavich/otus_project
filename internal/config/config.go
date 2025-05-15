package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Cache CacheConf
	HTTP  HTTPConf
}

type CacheConf struct {
	MaxItems int `env:"CACHE_ITEMS" env-default:"100"`
}

type HTTPConf struct {
	Host string `env:"HOST" env-default:"0.0.0.0"`
	Port int    `env:"PORT" env-default:"8081"`
}

func Load() *Config {
	var cfg Config
	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	return &cfg
}
