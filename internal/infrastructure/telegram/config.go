package telegram

import (
	"errors"
	"os"
)

type Config struct {
	Token string
	Debug bool
}

func LoadConfig() (*Config, error) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		return nil, errors.New("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	return &Config{
		Token: token,
		Debug: os.Getenv("DEBUG") == "true",
	}, nil
}
