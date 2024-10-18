package env

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from a .env file.
func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Loading sample env if available...")
		if err := godotenv.Load(".env.sample"); err != nil {
			log.Println("No .env.sample file found either.")
			return err
		}
	}
	return nil
}

// GetInterval gets the interval for sleeping between checks, with a default value.
func GetInterval() int {
	intervalStr := os.Getenv("INTERVAL")
	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		return 60 // Default to 60 seconds if not set
	}
	return interval
}
