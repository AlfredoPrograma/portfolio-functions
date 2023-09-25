package notion

import (
	"context"
	"reflect"

	"github.com/AlfredoPrograma/portfolio-functions/proposals/register/config"
	"github.com/AlfredoPrograma/portfolio-functions/proposals/register/errors"
	"github.com/jomei/notionapi"
)

type Client struct {
	databaseId string
	api        *notionapi.Client
}

func NewClient() Client {
	return Client{
		api:        notionapi.NewClient(notionapi.Token(config.Use().NOTION_API_KEY)),
		databaseId: config.Use().NOTION_DATABASE_ID,
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

func (c *Client) RegisterProposal(payload map[string]any) (*notionapi.Page, error) {
	data, err := validateProposal(payload)

	if err != nil {
		return nil, err
	}

	page, err := c.api.Page.Create(context.Background(), &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: notionapi.DatabaseID(c.databaseId),
		},
		Properties: notionapi.Properties{
			"Subject": notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{
						Text: &notionapi.Text{
							Content: data.Subject,
						},
					},
				},
			},
			"Email": notionapi.RichTextProperty{
				Type: "rich_text",
				RichText: []notionapi.RichText{
					{
						Text: &notionapi.Text{
							Content: data.Email,
						},
					},
				},
			},
			"Full Name": notionapi.RichTextProperty{
				Type: "rich_text",
				RichText: []notionapi.RichText{
					{
						Text: &notionapi.Text{
							Content: data.FullName,
						},
					},
				},
			},
			"Reviewed": notionapi.CheckboxProperty{
				Type:     "checkbox",
				Checkbox: false,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return page, nil
}
