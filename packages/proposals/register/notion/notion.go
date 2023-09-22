package notion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AlfredoPrograma/portfolio-functions/proposals/register/config"
	"github.com/AlfredoPrograma/portfolio-functions/proposals/register/errors"
)

type Client struct {
	apiKey     string
	baseUrl    string
	databaseId string
	version    string
	args       map[string]any
	httpClient *http.Client
}

func NewClient(args map[string]any) Client {
	return Client{
		apiKey:     config.Use().NOTION_API_KEY,
		baseUrl:    config.Use().NOTION_BASE_URL,
		databaseId: config.Use().NOTION_DATABASE_ID,
		version:    config.Use().NOTION_VERSION,
		httpClient: &http.Client{},
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

func (c *Client) RegisterProposal() (any, error) {
	payload, err := getPayload(c.args)

	if err != nil {
		return payload, err
	}

	notionPayload := map[string]any{
		"parent": map[string]any{
			"database_id": c.databaseId,
		},
		"properties": map[string]any{
			"Subject": map[string]any{
				"title": []map[string]any{
					{
						"text": map[string]string{
							"content": payload.Subject,
						},
					},
				},
			},
			"Email": map[string]any{
				"rich_text": []map[string]any{
					{
						"text": map[string]string{
							"content": payload.Email,
						},
					},
				},
			},
			"Full Name": map[string]any{
				"rich_text": []map[string]any{
					{
						"text": map[string]string{
							"content": payload.FullName,
						},
					},
				},
			},
			"Message": map[string]any{
				"rich_text": []map[string]any{
					{
						"text": map[string]string{
							"content": payload.Message,
						},
					},
				},
			},
		},
	}

	buffer := new(bytes.Buffer)
	err = json.NewEncoder(buffer).Encode(notionPayload)

	if err != nil {
		return payload, err
	}

	endpoint := fmt.Sprintf("%s/pages", c.baseUrl)

	req, err := http.NewRequest(http.MethodPost, endpoint, buffer)

	if err != nil {
		return payload, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Add("Notion-Version", c.version)
	req.Header.Add("Content-Type", "application/json")

	_, err = c.httpClient.Do(req)

	if err != nil {
		return payload, err
	}

	if err != nil {
		return payload, err
	}

	return map[string]any{
		"body": map[string]any{
			"message": "Success",
		},
	}, nil
}
