package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env            string     `yaml:"env" env-default:"prod"`
	TelegramToken  string     `yaml:"telegram_token" env:"TELEGRAM_TOKEN" env-required:"true"`
	MigrationsPath string     `yaml:"migrations_path" env-default:"./migrations"`
	StoragePath    string     `yaml:"storage_path" env-default:"./storage/storage.db"`
	GRPC           GRPC       `yaml:"grpc"`
	RateLimit      RateLimit  `yaml:"rate_limit"`
	Backoffice     Backoffice `yaml:"backoffice"`
}

type Backoffice struct {
	Retries         int           `yaml:"retries" env-default:"3"`
	MessageTimer    time.Duration `yaml:"message_timer" env-default:"5m"`
	RetriesTimeout  time.Duration `yaml:"retries_timeout" env-default:"5s"`
	ResponseTimeout time.Duration `yaml:"response_timeout" env-default:"15s"`
}

type RateLimit struct {
	FillPeriod  time.Duration `yaml:"fill_period"  env-default:"1s"`
	BucketLimit int           `yaml:"bucket_limit"  env-default:"10"`
}

type GRPC struct {
	Host    string        `yaml:"host" env-default:"localhost"`
	Port    string        `yaml:"port" env-default:"50051"`
	Timeout time.Duration `yaml:"timeout" env-default:"600s"`
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
