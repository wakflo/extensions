package actions

import (
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

// GetFolderOperation structure and methods
type getFolderProps struct {
	FolderID string `json:"folder-id"`
}

type GetFolderOperation struct{}

func (o *GetFolderOperation) Name() string {
	return "Get Folder"
}

func (o *GetFolderOperation) Description() string {
	return "Retrieves details of a specific ClickUp folder by ID."
}

func (o *GetFolderOperation) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (o *GetFolderOperation) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getFolderDocs,
	}
}

func (o *GetFolderOperation) Icon() *string {
	icon := "material-symbols:folder"
	return &icon
}

func (o *GetFolderOperation) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"workspace-id": shared.RegisterWorkSpaceInput("Workspaces", "select a workspace", true),
		"space-id":     shared.RegisterSpacesInput("Spaces", "select a space", true),
		"folder-id":    shared.RegisterFoldersInput("Folders", "select a folder", true),
	}
}

func (o *GetFolderOperation) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	accessToken := ctx.Auth.AccessToken
	input, err := sdk.InputToTypeSafely[getFolderProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}
	url := "/v2/folder/" + input.FolderID

	folder, _ := shared.GetData(accessToken, url)

	return folder, nil
}

func (o *GetFolderOperation) Auth() *sdk.Auth {
	return nil
}

func (o *GetFolderOperation) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":                "folder123",
		"name":              "Example Folder",
		"orderindex":        "1",
		"override_statuses": false,
		"hidden":            false,
		"space": map[string]any{
			"id":   "space123",
			"name": "Space Name",
		},
		"task_count": "15",
		"lists": []map[string]any{
			{
				"id":      "list123",
				"name":    "List 1",
				"content": "List description",
			},
			{
				"id":      "list456",
				"name":    "List 2",
				"content": "Another list description",
			},
		},
	}
}

func (o *GetFolderOperation) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetFolderOperation() sdk.Action {
	return &GetFolderOperation{}
}
