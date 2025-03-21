package actions

import (
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

// CreateFolderOperation structure and methods
type createFolderProps struct {
	SpaceID     string `json:"space-id"`
	WorkspaceID string `json:"workspace-id"`
	Name        string `json:"name"`
}

type CreateFolderOperation struct{}

func (o *CreateFolderOperation) Name() string {
	return "Create Folder"
}

func (o *CreateFolderOperation) Description() string {
	return "Creates a new folder in a specified ClickUp space."
}

func (o *CreateFolderOperation) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (o *CreateFolderOperation) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createFolderDocs,
	}
}

func (o *CreateFolderOperation) Icon() *string {
	icon := "material-symbols:create-new-folder"
	return &icon
}

func (o *CreateFolderOperation) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"workspace-id": shared.GetWorkSpaceInput("Workspaces", "select a workspace", true),
		"space-id":     shared.GetSpacesInput("Spaces", "select a space", true),
		"name": autoform.NewShortTextField().
			SetDisplayName("Folder Name").
			SetDescription("The name of the folder").
			SetRequired(true).
			Build(),
	}
}

func (o *CreateFolderOperation) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	accessToken := ctx.Auth.AccessToken
	input, err := sdk.InputToTypeSafely[createFolderProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}
	reqURL := "/v2/space/" + input.SpaceID + "/folder"

	folder, err := shared.CreateItem(accessToken, input.Name, reqURL)
	if err != nil {
		return nil, err
	}

	return folder, nil
}

func (o *CreateFolderOperation) Auth() *sdk.Auth {
	return nil
}

func (o *CreateFolderOperation) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":                "folder123",
		"name":              "New Folder",
		"orderindex":        "1",
		"override_statuses": false,
		"hidden":            false,
		"space": map[string]any{
			"id":   "space123",
			"name": "Space Name",
		},
		"task_count": "0",
	}
}

func (o *CreateFolderOperation) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateFolderOperation() sdk.Action {
	return &CreateFolderOperation{}
}
