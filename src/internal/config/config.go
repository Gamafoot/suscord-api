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
		Port    string        `yaml:"port" env-default:"8000"`
		Timeout time.Duration `yaml:"timeout" env-default:"5s"`
	} `yaml:"server"`

	Database struct {
		URL      string `yaml:"url" env-required:"true"`
		LogLevel string `yaml:"log_level"`
	} `yaml:"database"`

	Redis struct {
		Addr     string `yaml:"addr" env-required:"true"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db" env-default:"0"`
	} `yaml:"redis"`

	Hash struct {
		Salt string `yaml:"salt" env-required:"true"`
	} `yaml:"hash"`

	CORS struct {
		Origins        []string `yaml:"origins"`
		AllowedMethods []string `yaml:"allowed_methods"`
		AllowedHeaders []string `yaml:"allowed_headers"`
	} `yaml:"cors"`

	Media struct {
		AllowedMedia []string `yaml:"allowed_media" env-required:"true"`
		Url          string   `yaml:"url" env-default:"/media/"`
		Folder       string   `yaml:"folder" env-required:"true"`
		MaxSize      string   `yaml:"max_size" env-default:"500M"`
	} `yaml:"media"`

	Static struct {
		Url    string `yaml:"url" env-default:"/static/"`
		Folder string `yaml:"folder" env-required:"true"`
	} `yaml:"static"`
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
