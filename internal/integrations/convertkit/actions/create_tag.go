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

type createTagActionProps struct {
	TagName string `json:"tag_name"`
}

type CreateTagAction struct{}

// Metadata returns metadata about the action
func (a *CreateTagAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_tag",
		DisplayName:   "Create Tag",
		Description:   "Create a new tag in your ConvertKit account.",
		Type:          core.ActionTypeAction,
		Documentation: createTagDocs,
		Icon:          "mdi:tag",
		SampleOutput: map[string]any{
			"tag": map[string]any{
				"id":         "123456",
				"name":       "Example Tag",
				"created_at": "2024-03-15T10:30:00Z",
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CreateTagAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_tag", "Create Tag")

	form.TextField("tag_name", "Tag Name").
		Placeholder("Enter tag name").
		Required(true).
		HelpText("Name of the tag to create")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CreateTagAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CreateTagAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[createTagActionProps](ctx)
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
		"tag": map[string]string{
			"name": input.TagName,
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	response, err := shared.GetConvertKitClient("/tags", http.MethodPost, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewCreateTagAction() sdk.Action {
	return &CreateTagAction{}
}
