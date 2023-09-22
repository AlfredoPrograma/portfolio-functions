package main

import (
	"fmt"
	"net/http"
)

type Payload struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Subject  string `json:"subject"`
	Message  string `json:"message"`
}

type MissingFieldError struct {
	field string
}

func (e *MissingFieldError) Error() string {
	return fmt.Sprintf("Missing field: %s", e.field)
}

func getPayload(args map[string]any) (Payload, error) {
	payloadKeys := []string{"fullName", "email", "subject", "message"}

	for _, key := range payloadKeys {
		_, ok := args[key]

		if !ok {
			return Payload{}, &MissingFieldError{field: key}
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
	response := make(map[string]any)

	httpContext, ok := args["http"].(map[string]any)

	if !ok {
		response["body"] = map[string]any{
			"error": "no http context",
		}

		return response
	}

	if httpContext["method"] != http.MethodPost {
		response["body"] = map[string]any{
			"error": "method not allowed",
		}

		return response
	}

	payload, err := getPayload(args)

	if err != nil {
		response["body"] = map[string]any{
			"err": err.Error(),
		}

		return response
	}

	response["body"] = map[string]any{
		"data": payload,
	}

	return response
}
