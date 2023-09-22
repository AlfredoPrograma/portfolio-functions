package notion

import (
	"github.com/AlfredoPrograma/portfolio-functions/proposals/register/config"
	"github.com/AlfredoPrograma/portfolio-functions/proposals/register/errors"
)

type Client struct {
	apiKey     string
	baseUrl    string
	databaseId string
	version    string
	args       map[string]any
}

func NewClient(args map[string]any) Client {
	return Client{
		apiKey:     config.Use().NOTION_API_KEY,
		baseUrl:    config.Use().NOTION_BASE_URL,
		databaseId: config.Use().NOTION_DATABASE_ID,
		version:    config.Use().NOTION_VERSION,
		args:       args,
	}
}

type payload struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Subject  string `json:"subject"`
	Message  string `json:"message"`
}

func getPayload(args map[string]any) (payload, error) {
	payloadKeys := []string{"fullName", "email", "subject", "message"}

	for _, key := range payloadKeys {
		_, ok := args[key]

		if !ok {
			return payload{}, &errors.MissingKeyError{
				Context: "PAYLOAD",
				Field:   key,
			}
		}

		_, ok = args[key].(string)

		if !ok {
			return payload{}, &errors.InvalidFieldError{
				Field: key,
			}
		}
	}

	return payload{
		FullName: args["fullName"].(string),
		Email:    args["email"].(string),
		Subject:  args["subject"].(string),
		Message:  args["message"].(string),
	}, nil
}

func (c *Client) RegisterProposal() (payload, error) {
	return getPayload(c.args)
}
