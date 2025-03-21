package actions

import (
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getFolderlessListProps struct {
	SpaceID     string `json:"space-id"`
	WorkspaceID string `json:"workspace-id"`
}

type GetFolderlesslistOperation struct{}

func (o *GetFolderlesslistOperation) Name() string {
	return "Get Folderless Lists"
}

func (o *GetFolderlesslistOperation) Description() string {
	return "Retrieves all lists directly in a space that do not belong to any folder."
}

func (o *GetFolderlesslistOperation) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (o *GetFolderlesslistOperation) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getFolderlessListDocs,
	}
}

func (o *GetFolderlesslistOperation) Icon() *string {
	icon := "material-symbols:folder-off-outline"
	return &icon
}

func (o *GetFolderlesslistOperation) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"workspace-id": shared.GetWorkSpaceInput("Workspaces", "select a workspace", true),
		"space-id":     shared.GetSpacesInput("Spaces", "select a space", true),
	}
}

func (o *GetFolderlesslistOperation) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	accessToken := ctx.Auth.AccessToken
	input, err := sdk.InputToTypeSafely[getFolderlessListProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}
	folderlessList, _ := shared.GetSpace(accessToken, input.SpaceID)

	return folderlessList, nil
}

func (o *GetFolderlesslistOperation) Auth() *sdk.Auth {
	return nil
}

func (o *GetFolderlesslistOperation) SampleData() sdkcore.JSON {
	return map[string]any{
		"lists": []map[string]any{
			{
				"id":         "list123",
				"name":       "Folderless List 1",
				"content":    "First list",
				"orderindex": "1",
			},
			{
				"id":         "list456",
				"name":       "Folderless List 2",
				"content":    "Second list",
				"orderindex": "2",
			},
		},
	}
}

func (o *GetFolderlesslistOperation) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetFolderlesslistOperation() sdk.Action {
	return &GetFolderlesslistOperation{}
}
