package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/convertkit/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type tagSubscriberActionProps struct {
	TagID int    `json:"tag_id"`
	Email string `json:"email"`
}

type TagSubscriberAction struct{}

// Metadata returns metadata about the action
func (a *TagSubscriberAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "tag_subscriber",
		DisplayName:   "Tag Subscriber",
		Description:   "Apply a tag to a subscriber in your ConvertKit account.",
		Type:          core.ActionTypeAction,
		Documentation: tagSubscriberDocs,
		Icon:          "mdi:tag-plus",
		SampleOutput: map[string]any{
			"subscription": map[string]any{
				"id":                "12345",
				"state":             "inactive",
				"created_at":        "2024-03-15T10:30:00Z",
				"source":            nil,
				"referrer":          nil,
				"subscribable_id":   "789",
				"subscribable_type": "tag",
				"subscriber": map[string]any{
					"id": "54321",
				},
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *TagSubscriberAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("tag_subscriber", "Tag Subscriber")

	form.NumberField("tag_id", "Tag ID").
		Placeholder("Enter tag ID").
		Required(true).
		HelpText("ID of the tag to apply")

	form.TextField("email", "Email").
		Placeholder("Enter email address").
		Required(true).
		HelpText("Email address of the subscriber")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *TagSubscriberAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *TagSubscriberAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[tagSubscriberActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"api_secret": authCtx.Extra["api-secret"],
		"email":      input.Email,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	path := fmt.Sprintf("/tags/%d/subscribe", input.TagID)

	response, err := shared.GetConvertKitClient(path, http.MethodPost, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewTagSubscriberAction() sdk.Action {
	return &TagSubscriberAction{}
}
