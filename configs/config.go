package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenWeatherMapAPIKey string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf(".env not found in current directory, trying parent directory: %v", err)

		err = godotenv.Load("../.env")
		if err != nil {
			log.Printf("Error loading .env file: %v", err)
			return nil, err
		}
	}

	return &Config{
		OpenWeatherMapAPIKey: os.Getenv("OPENWEATHERMAP_API_KEY"),
	}, nil
}
