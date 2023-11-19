package config

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

var (
	cfg  Config
	once sync.Once
)

type HTTPServer struct {
	Address     string        `yaml:"address" env:"HTTP_SERVER_ADDRESS" env-default:"localhost:8082"`
	Timeout     time.Duration `yaml:"timeout" env:"HTTP_SERVER_TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"60s"`
}

type Goose struct {
	PrintStatus bool `yaml:"print_status"`
}

type Config struct {
	Env         string     `yaml:"env" env:"ENV" env-default:"local" env-required:"true"`
	StoragePath string     `yaml:"storage_path" env:"STORAGE_PATH" env-required:"true"`
	Goose       Goose      `yaml:"goose"`
	HttpServer  HTTPServer `yaml:"http_server"`
}

func mustLoad() {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		log.Fatal("[ERROR]: CONFIG_PATH is not set in PATH (env variable not set)")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("[ERROR]: config file %s does not exists", cfgPath)
	}

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		log.Fatalf("[ERROR]: read config completed with error: %v", err)
	}
}

func Instance() Config {
	once.Do(func() {
		mustLoad()
	})

	return cfg
}
