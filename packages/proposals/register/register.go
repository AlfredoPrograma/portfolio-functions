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
		return errors.NewErrorResponse(err.Error())
	}

	httpContext, ok := args["http"].(map[string]any)

	if !ok {
		return errors.NewErrorResponse("No http context")
	}

	if httpContext["method"] != http.MethodPost {
		return errors.NewErrorResponse(fmt.Sprintf("Method no allowed: %s", httpContext["method"]))
	}

	notion := notion.NewClient(args)
	payload, err := notion.RegisterProposal()

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
