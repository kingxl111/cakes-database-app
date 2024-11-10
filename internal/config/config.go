package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config type
type Config struct {
	Env        string     `yaml:"env" env-default:"local"`
	HTTPServer HTTPServer `yaml:"http_server"`
	DB         DB         `yaml:"db"`
}

// HTTPServer type
type HTTPServer struct {
	Address     string `yaml:"address" env-default:"localhost:8080"`
	Timeout     string `yaml:"timeout" env-default:"4s"`
	IdleTimeout string `yaml:"idle_timeout" env-default:"60s"`
}

// DB type
type DB struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DBName   string `yaml:"name"`
	SSLmode  string `yaml:"sslmode"`
}

const (
	user     = "DB_USER"
	password = "DB_PASSWORD"
	name     = "DB_NAME"
	host     = "DB_HOST"
	port     = "DB_PORT"
	sslmode  = "disable"
)

func MustLoad() *Config {

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if the file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist : %s", configPath) // Exit() - remember it!
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("can't read config: %s", err)
	}

	cfg.DB.DBName = os.Getenv(name)
	cfg.DB.Username = os.Getenv(user)
	cfg.DB.Password = os.Getenv(password)
	cfg.DB.Host = os.Getenv(host)
	cfg.DB.Port = os.Getenv(port)
	cfg.DB.SSLmode = sslmode

	return &cfg
}
