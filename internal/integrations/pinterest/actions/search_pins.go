package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/pinterest/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type searchPinsActionProps struct {
	Keyword           string `json:"keyword"`
	Bookmark          string `json:"bookmark,omitempty"`
	IsBusinessAccount string `json:"is_busness_account,omitempty"`
	AdAccountID       string `json:"ad_account_id,omitempty"`
}

type SearchPinsAction struct{}

// Metadata returns metadata about the action
func (a *SearchPinsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "search_pins",
		DisplayName:   "Search Pins",
		Description:   "Search pins within your Pinterest account (Note: Pinterest API v5 only searches your own pins, not all of Pinterest)",
		Type:          core.ActionTypeAction,
		Documentation: getPinDocs,
		Icon:          "mdi:magnify",
		SampleOutput: `{
			"items": [
				{
					"id": "123456789012345678",
					"created_at": "2024-01-15T10:30:00Z",
					"link": "https://www.example.com",
					"title": "My Pin Title",
					"description": "Pin description",
					"alt_text": null,
					"board_id": "987654321098765432",
					"board_section_id": null,
					"board_owner": {
						"id": "111222333444555666",
						"username": "myusername"
					},
					"media": {
						"media_type": "image",
						"images": {
							"600x": {
								"width": 600,
								"height": 900,
								"url": "https://i.pinimg.com/600x/..."
							},
							"original": {
								"width": 1200,
								"height": 1800,
								"url": "https://i.pinimg.com/originals/..."
							}
						}
					}
				}
			],
			"bookmark": "LT4xMjM0NTY3ODkwMTIzNDU2Nzg6NjU0MzIxfGU5ZDM2YjJmYzQ3YTI3MjM5ZGE5YzVhZWY3ZTI3YzFi"
		}`,
		Settings: core.ActionSettings{},
	}
}

func (a *SearchPinsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("search_pins", "Search Pins")

	form.TextField("keyword", "Search Keyword").
		Placeholder("Enter search term").
		Required(true).
		HelpText("Enter search term or phrase (e.g., 'summer recipes', 'home decor', 'fashion ideas')")

	form.SelectField("is_busness_account", "Are you using a business account?").
		Required(false).
		AddOption("yes", "Yes").
		AddOption("no", "No").
		HelpText("Confirm if you're using business account or a personal account")

	shared.GetAdAccountProps(form)

	form.TextField("bookmark", "Bookmark (Optional)").
		Placeholder("Leave empty for first page").
		Required(false).
		HelpText("Pagination cursor from previous search. Leave empty to get the first page of results.")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *SearchPinsAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *SearchPinsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[searchPinsActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Token.AccessToken == "" {
		return nil, errors.New("missing Pinterest auth token")
	}
	accessToken := authCtx.Token.AccessToken

	if input.Keyword == "" {
		return nil, errors.New("keyword is required")
	}

	// Search pins with optional bookmark for pagination and ad_account_id
	results, err := shared.SearchPins(accessToken, input.Keyword, input.Bookmark, input.AdAccountID)
	if err != nil {
		return nil, err
	}

	// Check if results is empty and provide helpful feedback
	if items, ok := results["items"].([]interface{}); ok && len(items) == 0 {
		return map[string]interface{}{
			"items":    []interface{}{},
			"message":  "No pins found. Note: Pinterest API v5 only searches within your own saved pins, not all of Pinterest.",
			"bookmark": nil,
		}, nil
	}

	return results, nil
}

func NewSearchPinsAction() sdk.Action {
	return &SearchPinsAction{}
}
