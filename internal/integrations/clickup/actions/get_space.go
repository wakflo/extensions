package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getSpaceProps struct {
	WorkspaceID string `json:"workspace-id"`
	SpaceID     string `json:"space-id"`
}

type GetSpaceOperation struct{}

// Metadata returns metadata about the action
func (o *GetSpaceOperation) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_space",
		DisplayName:   "Get Space",
		Description:   "Retrieves details of a specific ClickUp space",
		Type:          core.ActionTypeAction,
		Documentation: getSpaceDocs,
		Icon:          "material-symbols:space-dashboard",
		SampleOutput: map[string]any{
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
			"multiple_assignees": true,
			"features": map[string]any{
				"due_dates":     true,
				"time_tracking": true,
				"tags":          true,
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (o *GetSpaceOperation) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_space", "Get Space")

	shared.RegisterWorkSpaceInput(form, "Workspaces", "select a workspace", true)

	shared.RegisterSpacesInput(form, "Spaces", "select a space to retrieve", true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (o *GetSpaceOperation) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (o *GetSpaceOperation) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getSpaceProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	accessToken := authCtx.Token.AccessToken
	space, err := shared.GetSpace(accessToken, input.SpaceID)
	if err != nil {
		return nil, err
	}

	return space, nil
}

func NewGetSpaceOperation() sdk.Action {
	return &GetSpaceOperation{}
}
