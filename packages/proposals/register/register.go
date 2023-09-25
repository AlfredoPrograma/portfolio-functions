package main

import (
	"fmt"
	"net/http"

	"github.com/AlfredoPrograma/portfolio-functions/proposals/register/config"
	"github.com/AlfredoPrograma/portfolio-functions/proposals/register/errors"
	"github.com/AlfredoPrograma/portfolio-functions/proposals/register/notion"
)

func Main(args map[string]any) map[string]any {
	err := config.LoadEnv()

	if err != nil {
		return errors.NewErrorResponse("Error loading environment variables")
	}

	httpContext, ok := args["http"].(map[string]any)

	if !ok {
		return errors.NewErrorResponse("No http context")
	}

	if httpContext["method"] != http.MethodPost {
		return errors.NewErrorResponse(fmt.Sprintf("Method no allowed: %s", httpContext["method"]))
	}

	payload, ok := args["proposal"].(map[string]any)

	if !ok {
		return errors.NewErrorResponse("No payload given")
	}

	notion := notion.NewClient()
	page, err := notion.RegisterProposal(payload)

	if err != nil {
		return errors.NewErrorResponse("Error registering proposal")
	}

	return map[string]any{
		"body": map[string]any{
			"message": "Proposal registered",
			"id":      page.ID,
		},
	}
}
