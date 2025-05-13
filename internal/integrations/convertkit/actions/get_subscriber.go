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

type getSubscriberActionProps struct {
	SubscriberID int `json:"subscriber_id"`
}

type GetSubscriberAction struct{}

// Metadata returns metadata about the action
func (a *GetSubscriberAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_subscriber",
		DisplayName:   "Get Subscriber",
		Description:   "Retrieve detailed information about a specific subscriber by ID.",
		Type:          core.ActionTypeAction,
		Documentation: getSubscriberDocs,
		Icon:          "mdi:account-details",
		SampleOutput: map[string]any{
			"subscriber": map[string]any{
				"id":            "12345",
				"first_name":    "Jon",
				"email_address": "jon.snow@example.com",
				"state":         "active",
				"created_at":    "2023-05-15T10:30:00Z",
				"fields": map[string]any{
					"last_name": "Snow",
					"company":   "Night's Watch",
					"location":  "The Wall",
				},
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GetSubscriberAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_subscriber", "Get Subscriber")

	form.NumberField("subscriber_id", "Subscriber ID").
		Placeholder("Enter subscriber ID").
		Required(true).
		HelpText("ID of the subscriber to retrieve")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetSubscriberAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GetSubscriberAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[getSubscriberActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/subscribers/%d?api_secret=%s", input.SubscriberID, authCtx.Extra["api-secret"])

	response, err := shared.GetConvertKitClient(path, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewGetSubscriberAction() sdk.Action {
	return &GetSubscriberAction{}
}
