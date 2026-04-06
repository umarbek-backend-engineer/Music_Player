package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Api_Port string
}

func Load() *Config {
	err := godotenv.Load("./cmd/.env")
	if err != nil {
		log.Println(err)
		log.Fatal("Error in load configurations")
	}

	return &Config{
		Api_Port: os.Getenv("API_PORT"),
	}
}
