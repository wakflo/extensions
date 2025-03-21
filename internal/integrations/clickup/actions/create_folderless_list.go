package actions

import (
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createFolderlessListProps struct {
	SpaceID     string `json:"space-id"`
	WorkspaceID string `json:"workspace-id"`
	Name        string `json:"name"`
}

type CreateFolderlessListOperation struct{}

func (o *CreateFolderlessListOperation) Name() string {
	return "Create Folderless List"
}

func (o *CreateFolderlessListOperation) Description() string {
	return "Creates a new list directly in a space without a parent folder."
}

func (o *CreateFolderlessListOperation) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (o *CreateFolderlessListOperation) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createFolderlessListDocs,
	}
}

func (o *CreateFolderlessListOperation) Icon() *string {
	icon := "material-symbols:add-box-outline"
	return &icon
}

func (o *CreateFolderlessListOperation) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"workspace-id": shared.GetWorkSpaceInput("Workspaces", "select a workspace", true),
		"space-id":     shared.GetSpacesInput("Spaces", "select a space", true),
		"name": autoform.NewShortTextField().
			SetDisplayName("List Name").
			SetDescription("The name of the list").
			SetRequired(true).
			Build(),
	}
}

func (o *CreateFolderlessListOperation) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	accessToken := ctx.Auth.AccessToken
	input, err := sdk.InputToTypeSafely[createFolderlessListProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}
	reqURL := "/v2/space/" + input.SpaceID + "/list"

	response, err := shared.CreateItem(accessToken, input.Name, reqURL)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (o *CreateFolderlessListOperation) Auth() *sdk.Auth {
	return nil
}

func (o *CreateFolderlessListOperation) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":         "list123",
		"name":       "New Folderless List",
		"content":    "List description",
		"orderindex": 1,
		"status": map[string]any{
			"status": "Open",
			"color":  "#d3d3d3",
		},
	}
}

func (o *CreateFolderlessListOperation) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateFolderlessListOperation() sdk.Action {
	return &CreateFolderlessListOperation{}
}
