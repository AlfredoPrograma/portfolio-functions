package main

import (
	"fmt"
	"net/http"

	"github.com/AlfredoPrograma/portfolio-functions/proposals/register/config"
	"github.com/AlfredoPrograma/portfolio-functions/proposals/register/errors"
)

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
	err := config.LoadEnv()

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
			"environment": config.Use(),
		},
	}
}
