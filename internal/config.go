package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/lib/pq"
)

type Config struct {
	Env        string `yaml:"env"`
	Storage    string `yaml:"storage" env-required:"true"`
	HTTP_serer `yaml:"http_server"`
}
type HTTP_serer struct {
	Address      string        `yaml:"address" env-default:"localhost:8080"`
	Timeout      time.Duration `yaml:"timeout" env-default:"5s"`
	Idle_tymeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func FillConfig() Config {
	cfg := Config{}
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("Config_path is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist by path:%s", configPath)
	}
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config:%s", err)
	}
	return cfg
}
