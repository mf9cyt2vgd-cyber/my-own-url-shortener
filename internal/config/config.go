package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	User       User       `yaml:"user"`
	HttpServer HttpServer `yaml:"http_server"`
	Postgresql Postgresql `yaml:"postgresql"`
}
type User struct {
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
}
type HttpServer struct {
	Address     string        `yaml:"address"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
	Timeout     time.Duration `yaml:"timeout"`
}
type Postgresql struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	Sslmode  string `yaml:"sslmode"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG PATH is not set")
	}

	//check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("config file %s does not exist", configPath)
	}

	return &cfg
}
