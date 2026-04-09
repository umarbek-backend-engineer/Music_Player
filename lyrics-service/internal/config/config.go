package config

import (
	"lyrics-service/pkg/utils"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Api_Port         string
	Api_Host         string
	NetWork_Protocol string
	DB_Host          string
	DB_Port          string
	DB_User          string
	DB_Password      string
	DB_Driver        string
	DB_Name          string
}

func Load() *Config {

	err := godotenv.Load("./cmd/.env")
	if err != nil {
		err = utils.MapError(err)
	}

	return &Config{
		Api_Port:         os.Getenv("API_PORT"),
		Api_Host:         os.Getenv("API_HOST"),
		NetWork_Protocol: os.Getenv("NETWORK_PROTOCOL"),
		DB_Host:          os.Getenv("DB_HOST"),
		DB_Port:          os.Getenv("DB_PORT"),
		DB_User:          os.Getenv("DB_USER"),
		DB_Password:      os.Getenv("DB_PASSWORD"),
		DB_Driver:        os.Getenv("DB_DRIVER"),
		DB_Name:          os.Getenv("DB_NAME"),
	}
}
