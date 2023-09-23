package notion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

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

type proposal struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Subject  string `json:"subject"`
	Message  string `json:"message"`
}

func validateProposal(payload map[string]any) (proposal, error) {
	payloadKeys := []string{"fullName", "email", "subject", "message"}

	for _, key := range payloadKeys {
		_, ok := payload[key]

		// Check if some key is missing
		if !ok {
			return proposal{}, &errors.MissingKeyError{
				Context: "PAYLOAD",
				Field:   key,
			}
		}

		fieldType := reflect.TypeOf(payload[key]).Kind()

		// Check if some key is not a string
		if fieldType != reflect.String {
			return proposal{}, &errors.InvalidFieldError{
				Field: key,
			}
		}

		if len(payload[key].(string)) == 0 {
			return proposal{}, &errors.EmptyFieldError{
				Field: key,
			}
		}
	}

	return proposal{
		FullName: payload["fullName"].(string),
		Email:    payload["email"].(string),
		Subject:  payload["subject"].(string),
		Message:  payload["message"].(string),
	}, nil
}

func (c *Client) RegisterProposal(payload map[string]any) (any, error) {
	data, err := validateProposal(payload)

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
							"content": data.Subject,
						},
					},
				},
			},
			"Email": map[string]any{
				"rich_text": []map[string]any{
					{
						"text": map[string]string{
							"content": data.Email,
						},
					},
				},
			},
			"Full Name": map[string]any{
				"rich_text": []map[string]any{
					{
						"text": map[string]string{
							"content": data.FullName,
						},
					},
				},
			},
			"Message": map[string]any{
				"rich_text": []map[string]any{
					{
						"text": map[string]string{
							"content": data.Message,
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
