package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB struct {
		Host     string `env:"DB_HOST" required:"true"`
		Port     string `env:"DB_PORT,required" default:"5432"`
		Password string `env:"DB_PASSWORD"`
		Name     string `env:"DB_NAME"`
	} `env:"DB"`
	Server struct {
		Addr string `env:"SERVER_ADDR"`
		Port string `env:"SERVER_PORT"`
	} `env:"SERVER"`
	Migrate struct {
		Path string `env:"MIGRATE_PATH"`
	}
}

func NewConfig() *Config {
	var cfg Config

	if err := godotenv.Load("config.env"); err != nil {
		log.Fatalf("Couldn't load configuration .env")
	}
	cfg.DB.Host = os.Getenv("DB_HOST")
	cfg.DB.Name = os.Getenv("DB_NAME")
	cfg.DB.Port = os.Getenv("DB_PORT")
	cfg.DB.Password = os.Getenv("DB_PASSWORD")
	cfg.Server.Addr = os.Getenv("SERVER_ADDR")
	cfg.Server.Port = os.Getenv("SERVER_PORT")
	cfg.Migrate.Path = os.Getenv("MIGRATE_PATH")
	return &cfg
}
