package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// configuration items.
type Config struct {
	API_Port        string
	DB_Host         string
	DB_Port         string
	DB_User         string
	DB_Password     string
	DB_Name         string
	DB_Driver       string
	StoragePath     string
	NetworkProtocol string
	Rabbit_Port     string
	Rabbit_User     string
	Rabbit_Password string
	Rabbit_Host     string
}

func Load() *Config {
	err := godotenv.Load("./cmd/.env")
	if err != nil {
		log.Fatal("Error in loading config file \n", err)
	}

	return &Config{
		API_Port:        os.Getenv("API_PORT"),
		DB_Host:         os.Getenv("DB_HOST"),
		DB_Port:         os.Getenv("DB_PORT"),
		DB_User:         os.Getenv("DB_USER"),
		DB_Password:     os.Getenv("DB_PASSWORD"),
		DB_Name:         os.Getenv("DB_NAME"),
		DB_Driver:       os.Getenv("DB_DRIVER"),
		StoragePath:     os.Getenv("STORAGEPATH"),
		NetworkProtocol: os.Getenv("NETWORK_PROTOCOL"),
		Rabbit_Port:     os.Getenv("RABBIT_PORT"),
		Rabbit_User:     os.Getenv("RABBIT_USER"),
		Rabbit_Password: os.Getenv("RABBIT_PASSWORD"),
		Rabbit_Host:     os.Getenv("RABBIT_HOST"),
	}

}
