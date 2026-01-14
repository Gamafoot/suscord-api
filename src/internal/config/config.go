package config

import (
	"flag"
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server struct {
		Port    string        `yaml:"port" env:"SERVER_PORT" env-default:"8000"`
		Timeout time.Duration `yaml:"timeout" env:"SERVER_TIMEOUT" env-default:"10s"`
	} `yaml:"server"`

	Broker struct {
		Addr     string `yaml:"addr" env:"BROKER_ADDR" env-required:"true"`
		PoolSize int    `yaml:"pool_size" env:"BROKER_POOL_SIZE" env-default:"3"`
	} `yaml:"broker"`

	Database struct {
		URL      string `yaml:"url" env:"DB_URL" env-required:"true"`
		LogLevel string `yaml:"log_level" env:"DB_LOG_LEVEL"`
	} `yaml:"database"`

	Redis struct {
		Addr     string `yaml:"addr" env:"REDIS_ADDR" env-required:"true"`
		Password string `yaml:"password" env:"REDIS_PASS"`
		DB       int    `yaml:"db" env:"REDIS_DB" env-default:"0"`
	} `yaml:"redis"`

	Hash struct {
		Salt string `yaml:"salt" env:"HASH_SALT" env-required:"true"`
	} `yaml:"hash"`

	CORS struct {
		Origins        []string `yaml:"origins"`
		AllowedMethods []string `yaml:"allowed_methods"`
		AllowedHeaders []string `yaml:"allowed_headers"`
	} `yaml:"cors"`

	Media struct {
		Url          string   `yaml:"url" env:"MEDIA_URL" env-default:"/media/"`
		Folder       string   `yaml:"folder" env:"MEDIA_FOLDER" env-required:"true"`
		AllowedMedia []string `yaml:"allowed_media" env-required:"true"`
		MaxSize      string   `yaml:"max_size" env:"MEDIA_MAX_SIZE" env-default:"500M"`
	} `yaml:"media"`

	Static struct {
		URL    string `yaml:"url" env:"STATIC_URL" env-default:"/static/"`
		Folder string `yaml:"folder" env:"STATIC_FOLDER" env-required:"true"`
	} `yaml:"static"`

	Logger struct {
		Level  string `yaml:"level" env:"LOGGER_LEVEL" env-default:"info"`
		Folder string `yaml:"folder" env:"LOGGER_FOLDER" env-default:"assets/log"`
	} `yaml:"logger"`
}

var cfg *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		cfg = &Config{}

		path := getConfigPath()

		if err := cleanenv.ReadConfig(path, cfg); err != nil {
			panic(err)
		}
	})

	return cfg
}

func getConfigPath() string {
	var path string
	flag.StringVar(&path, "config", "../config/config.yaml", "set config file")

	envPath := os.Getenv("CONFIG_PATH")

	if len(envPath) > 0 {
		path = envPath
	}

	return path
}
