package actions

import (
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/convertkit/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type ListTagsAction struct{}

// Metadata returns metadata about the action
func (a *ListTagsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_tags",
		DisplayName:   "List Tags",
		Description:   "Retrieve a list of all tags in your ConvertKit account.",
		Type:          core.ActionTypeAction,
		Documentation: listTagsDocs,
		Icon:          "mdi:tag-multiple",
		SampleOutput: map[string]any{
			"tags": []map[string]any{
				{
					"id":         "1001",
					"name":       "Newsletter",
					"created_at": "2023-05-15T10:30:00Z",
				},
				{
					"id":         "1002",
					"name":       "Product Updates",
					"created_at": "2023-06-20T14:45:00Z",
				},
				{
					"id":         "1003",
					"name":       "New Customers",
					"created_at": "2023-07-12T09:15:00Z",
				},
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ListTagsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_tags", "List Tags")

	form.NumberField("limit", "Limit").
		Placeholder("Enter limit").
		Required(false).
		DefaultValue(50).
		HelpText("Maximum number of subscribers to retrieve (default: 50)")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ListTagsAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ListTagsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	// Construct the path with API key
	path := "/tags?api_key=" + authCtx.Extra["api-key"]

	// Make the API request
	response, err := shared.GetConvertKitClient(path, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	// Return the response
	return response, nil
}

func NewListTagsAction() sdk.Action {
	return &ListTagsAction{}
}
