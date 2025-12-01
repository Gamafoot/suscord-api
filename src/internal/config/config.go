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
	} `yaml:"listen"`

	Database struct {
		URL string `env:"url" env-required:"true"`
	} `yaml:"database"`

	Hash struct {
		Salt string `yaml:"salt" env-required:"true"`
	} `yaml:"hash"`

	CORS struct {
		Origins []string `yaml:"origins"`
	} `yaml:"cors"`

	Media struct {
		MaxFileSize       int64    `yaml:"max_file_size" env-default:"734003200"`
		AllowedExtentions []string `yaml:"allowed_types" env-required:"true"`
		RootFolder        string   `yaml:"root_folder" env-required:"true"`
		RootUrl           string   `yaml:"root_url" env-required:"true"`
	} `yaml:"media"`

	Static struct {
		RootFolder string `yaml:"root_folder" env-required:"true"`
		RootUrl    string `yaml:"root_url" env-required:"true"`
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
