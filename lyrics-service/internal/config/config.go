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
	}
}
