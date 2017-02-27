package config

import (
	"os"

	"github.com/subosito/gotenv"
)

// LoadEnv loads a json configuration file
func LoadEnv(path string) error {
	if path == "" {
		path = ".env"
	}

	return gotenv.Load(path)
}

// IsProductionEnvironment returns true if environment is production
func IsProductionEnvironment() bool {
	environment := os.Getenv("ENVIRONMENT")
	if environment == "production" {
		return true
	}
	return false
}
