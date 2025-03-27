package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Env            string `yaml:"env" env-default:"prod"`
	TelegramToken  string `yaml:"telegram_token" env:"TELEGRAM_TOKEN" env-required:"true"`
	MigrationsPath string `yaml:"migrations_path" env-default:"./migrations"`
	GRPC           GRPC   `yaml:"grpc"`
}

type GRPC struct {
	Host    string `yaml:"host" env-default:"localhost"`
	Port    string `yaml:"port" env-default:"50051"`
	Timeout string `yaml:"timeout" env-default:"600s"`
}

func MustLoad() *Config {
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}
	if configPath == "" {
		panic("config file path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist : " + err.Error())
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("can't read config : " + err.Error())
	}

	return &cfg
}
