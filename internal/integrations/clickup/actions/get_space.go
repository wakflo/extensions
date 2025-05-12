package actions

import (
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getSpaceProps struct {
	WorkspaceID string `json:"workspace-id"`
	SpaceID     string `json:"space-id"`
}

type GetSpaceOperation struct{}

func (o *GetSpaceOperation) Name() string {
	return "Get Space"
}

func (o *GetSpaceOperation) Description() string {
	return "Retrieves details of a specific ClickUp space"
}

func (o *GetSpaceOperation) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (o *GetSpaceOperation) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getSpaceDocs,
	}
}

func (o *GetSpaceOperation) Icon() *string {
	icon := "material-symbols:space-dashboard"
	return &icon
}

func (o *GetSpaceOperation) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"workspace-id": shared.RegisterWorkSpaceInput("Workspaces", "select a workspace", true),
		"space-id":     shared.RegisterSpacesInput("Spaces", "select a space to retrieve", true),
	}
}

func (o *GetSpaceOperation) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	accessToken := ctx.Auth.AccessToken

	input, err := sdk.InputToTypeSafely[getSpaceProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	space, err := shared.GetSpace(accessToken, input.SpaceID)
	if err != nil {
		return nil, err
	}

	return space, nil
}

func (o *GetSpaceOperation) Auth() *sdk.Auth {
	return nil
}

func (o *GetSpaceOperation) SampleData() sdkcore.JSON {
	return map[string]any{
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
	}
}

func (o *GetSpaceOperation) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetSpaceOperation() sdk.Action {
	return &GetSpaceOperation{}
}
