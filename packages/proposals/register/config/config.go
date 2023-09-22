package config

import (
	"os"

	"github.com/AlfredoPrograma/portfolio-functions/proposals/register/errors"
)

type Env struct {
	NOTION_API_KEY     string
	NOTION_DATABASE_ID string
	NOTION_BASE_URL    string
	NOTION_VERSION     string
}

var env = Env{}

func LoadEnv() error {
	envKeys := []string{"NOTION_API_KEY", "NOTION_DATABASE_ID", "NOTION_BASE_URL", "NOTION_VERSION"}

	for _, key := range envKeys {
		value, ok := os.LookupEnv(key)

		if !ok {
			return &errors.MissingKeyError{
				Context: "ENV",
				Field:   key,
			}
		}

		switch key {
		case "NOTION_API_KEY":
			env.NOTION_API_KEY = value
		case "NOTION_DATABASE_ID":
			env.NOTION_DATABASE_ID = value
		case "NOTION_BASE_URL":
			env.NOTION_BASE_URL = value
		case "NOTION_VERSION":
			env.NOTION_VERSION = value
		}
	}

	return nil
}

func Use() Env {
	return env
}
