package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC        gRPCConfig    `yaml:"grpc"`
	HTTP        HTTPServer    `yaml:"http"`
}

type gRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type HTTPServer struct {
	Port int `yaml:"port"`
}

type Secret struct {
	UserKey string `env:"USER_KEY" env-required:"true"`
}

func MustLoad() (*Config, *Secret) {
	configPath, envPath := fetchPath()

	if configPath == "" {
		log.Fatal("Config path is empty")
	}
	if envPath == "" {
		log.Fatal("Env path is empty")
	}

	scr := &Secret{}
	if err := cleanenv.ReadEnv(scr); err != nil {
		log.Fatalf("failed to get user key env")
	}

	return MustLoadConfigPath(configPath), scr
}

func MustLoadConfigPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file %s does not exist", configPath)
	}

	cfg := Config{}
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Println(configPath)
		log.Fatalf("Cannot read config: %s", err)
	}

	return &cfg
}

func fetchPath() (string, string) {
	var configPath, envPath string

	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.StringVar(&envPath, "env", "", "path to env file")
	flag.Parse()

	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Env file %s does not exist", envPath)
	}

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	return configPath, envPath
}
