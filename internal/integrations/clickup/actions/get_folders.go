package actions

import (
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getFoldersProps struct {
	SpaceID     string `json:"space-id"`
	WorkspaceID string `json:"workspace-id"`
}

type GetFoldersOperation struct{}

func (o *GetFoldersOperation) Name() string {
	return "Get Folders"
}

func (o *GetFoldersOperation) Description() string {
	return "Retrieves all folders within a specified ClickUp space."
}

func (o *GetFoldersOperation) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (o *GetFoldersOperation) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getFoldersDocs,
	}
}

func (o *GetFoldersOperation) Icon() *string {
	icon := "material-symbols:folder-copy"
	return &icon
}

func (o *GetFoldersOperation) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"workspace-id": shared.GetWorkSpaceInput("Workspaces", "select a workspace", true),
		"space-id":     shared.GetSpacesInput("Spaces", "select a space", true),
	}
}

func (o *GetFoldersOperation) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	accessToken := ctx.Auth.AccessToken

	input, err := sdk.InputToTypeSafely[getFoldersProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}
	url := "/v2/space/" + input.SpaceID + "/folder"

	folders, err := shared.GetData(accessToken, url)
	if err != nil {
		return nil, err
	}

	return folders, nil
}

func (o *GetFoldersOperation) Auth() *sdk.Auth {
	return nil
}

func (o *GetFoldersOperation) SampleData() sdkcore.JSON {
	return map[string]any{
		"folders": []map[string]any{
			{
				"id":                "folder123",
				"name":              "Folder 1",
				"orderindex":        "1",
				"override_statuses": false,
				"hidden":            false,
				"task_count":        "8",
			},
			{
				"id":                "folder456",
				"name":              "Folder 2",
				"orderindex":        "2",
				"override_statuses": false,
				"hidden":            false,
				"task_count":        "12",
			},
		},
	}
}

func (o *GetFoldersOperation) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetFoldersOperation() sdk.Action {
	return &GetFoldersOperation{}
}
