package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Api_port     string
	Api_protocol string
	DB_port      string
	DB_host      string
	DB_user      string
	DB_password  string
	DB_name      string
}

// this fucntion loads all existing variables from .env file to condig struct
func Load() *Config {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		Api_port:     os.Getenv("API_PORT"),
		Api_protocol: os.Getenv("API_PROTOCOL"),
		DB_port:      os.Getenv("DB_PORT"),
		DB_host:      os.Getenv("DB_HOST"),
		DB_user:      os.Getenv("DB_USER"),
		DB_password:  os.Getenv("DB_PASSWORD"),
		DB_name:      os.Getenv("DB_NAME"),
	}
}
