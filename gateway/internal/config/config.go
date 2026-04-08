package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Api_Port        string
	Api_Host        string
	Rabbit_Port     string
	Rabbit_User     string
	Rabbit_Password string
	Rabbit_Host     string
}

func Load() *Config {
	err := godotenv.Load("./cmd/.env")
	if err != nil {
		log.Println(err)
		log.Fatal("Error in load configurations")
	}

	return &Config{
		Api_Port:        os.Getenv("API_PORT"),
		Api_Host:        os.Getenv("API_HOST"),
		Rabbit_Port:     os.Getenv("RABBIT_PORT"),
		Rabbit_User:     os.Getenv("RABBIT_USER"),
		Rabbit_Password: os.Getenv("RABBIT_PASSWORD"),
		Rabbit_Host:     os.Getenv("RABBIT_HOST"),
	}
}
