package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPS        gRPCConfig    `yaml:"grpc"`
}

type gRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"Timeout"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		log.Fatal("Config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file %s does not exist", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot read config: %s", err)
	}

	return &cfg
}

func fetchConfigPath() string {
	var configPath string

	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	return configPath
}
