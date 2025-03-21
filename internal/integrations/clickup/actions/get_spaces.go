package actions

import (
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getSpacesProps struct {
	TeamID string `json:"team-id"`
}

type GetSpacesOperation struct{}

func (o *GetSpacesOperation) Name() string {
	return "Get Spaces"
}

func (o *GetSpacesOperation) Description() string {
	return "Retrieves all spaces in a ClickUp workspace."
}

func (o *GetSpacesOperation) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (o *GetSpacesOperation) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getSpacesDocs,
	}
}

func (o *GetSpacesOperation) Icon() *string {
	icon := "material-symbols:space-dashboard-outline"
	return &icon
}

func (o *GetSpacesOperation) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"team-id": shared.GetWorkSpaceInput("Workspaces", "select a workspace", true),
	}
}

func (o *GetSpacesOperation) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	accessToken := ctx.Auth.AccessToken
	input, err := sdk.InputToTypeSafely[getSpacesProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	spaces, _ := shared.GetAllSpaces(accessToken, input.TeamID)

	return spaces, nil
}

func (o *GetSpacesOperation) Auth() *sdk.Auth {
	return nil
}

func (o *GetSpacesOperation) SampleData() sdkcore.JSON {
	return map[string]any{
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
	}
}

func (o *GetSpacesOperation) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetSpacesOperation() sdk.Action {
	return &GetSpacesOperation{}
}
