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
	S3         S3         `yaml:"s3_config"`
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

type S3 struct {
	Endpoint  string
	Bucket    string
	AccessKey string
	SecretKey string
	Region    string
	PublicUrl string
}

const (
	user     = "DB_USER"
	password = "DB_PASSWORD"
	name     = "DB_NAME"
	host     = "DB_HOST"
	port     = "DB_PORT"
	sslmode  = "disable"

	s3Endpoint  = "S3_ENDPOINT"
	s3Bucket    = "S3_BUCKET"
	s3AccessKey = "S3_ACCESS_KEY"
	s3SecretKey = "S3_SECRET_KEY"
	s3Region    = "S3_REGION"
	s3PublicURL = "PUBLIC_S3_URL"
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

	cfg.S3.Endpoint = os.Getenv(s3Endpoint)
	cfg.S3.AccessKey = os.Getenv(s3AccessKey)
	cfg.S3.SecretKey = os.Getenv(s3SecretKey)
	cfg.S3.Region = os.Getenv(s3Region)
	cfg.S3.Bucket = os.Getenv(s3Bucket)
	cfg.S3.PublicUrl = os.Getenv(s3PublicURL)

	return &cfg
}
