package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/AlfredoPrograma/portfolio-functions/proposals/register/errors"
)

type Env struct {
	NOTION_API_KEY     string
	NOTION_DATABASE_ID string
	NOTION_BASE_URL    string
	NOTION_VERSION     string
}

func loadEnv() (Env, error) {
	env := &Env{}
	envKeys := []string{"NOTION_API_KEY", "NOTION_DATABASE_ID", "NOTION_BASE_URL", "NOTION_VERSION"}

	for _, key := range envKeys {
		value, ok := os.LookupEnv(key)

		if !ok {
			return Env{}, &errors.MissingKeyError{
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

	return *env, nil
}

type Payload struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Subject  string `json:"subject"`
	Message  string `json:"message"`
}

func getPayload(args map[string]any) (Payload, error) {
	payloadKeys := []string{"fullName", "email", "subject", "message"}

	for _, key := range payloadKeys {
		_, ok := args[key]

		if !ok {
			return Payload{}, &errors.MissingKeyError{
				Context: "PAYLOAD",
				Field:   key,
			}
		}

		_, ok = args[key].(string)

		if !ok {
			return Payload{}, &errors.InvalidFieldError{
				Field: key,
			}
		}
	}

	return Payload{
		FullName: args["fullName"].(string),
		Email:    args["email"].(string),
		Subject:  args["subject"].(string),
		Message:  args["message"].(string),
	}, nil
}

func Main(args map[string]any) map[string]any {
	env, err := loadEnv()

	if err != nil {
		return errors.NewErrorResponse(err.Error())
	}

	httpContext, ok := args["http"].(map[string]any)

	if !ok {
		return errors.NewErrorResponse("No http context")
	}

	if httpContext["method"] != http.MethodPost {
		return errors.NewErrorResponse(fmt.Sprintf("Method no allowed: %s", httpContext["method"]))
	}

	payload, err := getPayload(args)

	if err != nil {
		return errors.NewErrorResponse(err.Error())
	}

	return map[string]any{
		"body": map[string]any{
			"payload":     payload,
			"environment": env,
		},
	}
}
