package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getSpacesProps struct {
	TeamID string `json:"workspace-id"` // Changed to match kebab-case convention
}

type GetSpacesAction struct{}

// Metadata returns metadata about the action
func (a *GetSpacesAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_spaces",
		DisplayName:   "Get Spaces",
		Description:   "Retrieves all spaces in a ClickUp workspace.",
		Type:          core.ActionTypeAction,
		Documentation: getSpacesDocs,
		SampleOutput: map[string]any{
			"spaces": []map[string]any{
				{
					"id":      "123456",
					"name":    "Marketing Space",
					"private": false,
					"statuses": []map[string]any{
						{
							"id":     "st123",
							"status": "Open",
							"color":  "#d3d3d3",
						},
					},
				},
				{
					"id":      "789012",
					"name":    "Development Space",
					"private": true,
					"statuses": []map[string]any{
						{
							"id":     "st456",
							"status": "In Progress",
							"color":  "#4286f4",
						},
					},
				},
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GetSpacesAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_spaces", "Get Spaces")

	shared.RegisterWorkSpaceInput(form, "Workspace", "select a workspace", true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetSpacesAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GetSpacesAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[getSpacesProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the token source from the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	spaces, err := shared.GetAllSpaces(authCtx.Token.AccessToken, input.TeamID)
	if err != nil {
		return nil, err
	}

	return spaces, nil
}

func NewGetSpacesAction() sdk.Action {
	return &GetSpacesAction{}
}
