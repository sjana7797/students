package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Address string `yaml:"address"`
}

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HttpServer  `yaml:"http_server"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "path to the configuration file")

		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("Config path not set")
		}
	}

	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		log.Fatalf("COnfig file does not exist: %s", configPath)
	}

	var cfg Config

	err = cleanenv.ReadConfig(configPath, &cfg)

	if err != nil {
		log.Fatal("Can not read config file")
	}

	return &cfg
}
