package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env     string `yaml:"env" env:"ENV" env-default:"prod" env-required:"true"`
	Storage `yaml:"storage"`
	Server  `yaml:"http_server"`
}

type Storage struct {
	Host         string `yaml:"host" env:"DB-HOST" env-required:"true"`
	Port         int    `yaml:"port" env:"DB-PORT" env-required:"true"`
	User         string `yaml:"user" env:"DB-USER" env-required:"true"`
	Password     string `yaml:"password" env:"DB-PASSWORD" env-required:"true"`
	DatabaseName string `yaml:"dbname" env:"DB-NAME" env-required:"true"`
}

type Server struct {
	Address     string        `yaml:"address" env-default:"localhost:8881"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("cannot read config: %s", err)
	}

	return &cfg
}
