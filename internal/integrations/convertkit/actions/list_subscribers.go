package actions

import (
	"fmt"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/convertkit/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type listSubscribersActionProps struct {
	Page     int  `json:"page"`
	PageSize int  `json:"page-size"`
	FromDate *int `json:"from-date"`
}

type ListSubscribersAction struct{}

// Metadata returns metadata about the action
func (a *ListSubscribersAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_subscribers",
		DisplayName:   "List Subscribers",
		Description:   "Retrieve a list of subscribers from your ConvertKit account with their details and tags.",
		Type:          core.ActionTypeAction,
		Documentation: listSubscribersDocs,
		Icon:          "mdi:account-group",
		SampleOutput: map[string]any{
			"subscribers": []map[string]any{
				{
					"id":            "123456",
					"first_name":    "Jane",
					"email_address": "jane@example.com",
					"state":         "active",
					"created_at":    "2023-01-15T10:30:00Z",
					"fields": map[string]string{
						"company": "Acme Inc",
					},
				},
				{
					"id":            "789012",
					"first_name":    "John",
					"email_address": "john@example.com",
					"state":         "active",
					"created_at":    "2023-01-16T14:20:00Z",
					"fields": map[string]string{
						"company": "XYZ Corp",
					},
				},
			},
			"total_subscribers": "156",
			"page":              "1",
			"page_size":         "50",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ListSubscribersAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_subscribers", "List Subscribers")

	form.NumberField("page", "Page").
		Placeholder("Enter page number").
		Required(false).
		HelpText("The page number to retrieve (pagination)")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ListSubscribersAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ListSubscribersAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	_, err := sdk.InputToTypeSafely[listSubscribersActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	path := "/subscribers?api_secret=" + authCtx.Extra["api-secret"]

	response, err := shared.GetConvertKitClient(path, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	responseMap, ok := response.(map[string]interface{})
	if !ok {
		return nil, err
	}

	subscribers, ok := responseMap["subscribers"]
	if !ok {
		return nil, fmt.Errorf("failed to extract subscribers from response")
	}

	return subscribers, nil
}

func NewListSubscribersAction() sdk.Action {
	return &ListSubscribersAction{}
}
