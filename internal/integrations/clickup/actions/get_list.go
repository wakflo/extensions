package actions

import (
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getListProps struct {
	ListID string `json:"list-id"`
}

type GetListOperation struct{}

func (o *GetListOperation) Name() string {
	return "Get List"
}

func (o *GetListOperation) Description() string {
	return "Retrieves details of a specific ClickUp list by ID."
}

func (o *GetListOperation) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (o *GetListOperation) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getListDocs,
	}
}

func (o *GetListOperation) Icon() *string {
	icon := "material-symbols:format-list-bulleted"
	return &icon
}

func (o *GetListOperation) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"workspace-id": shared.GetWorkSpaceInput("Workspaces", "select a workspace", true),
		"space-id":     shared.GetSpacesInput("Spaces", "select a space", true),
		"folder-id":    shared.GetFoldersInput("Folders", "select a folder", true),
		"list-id":      shared.GetListsInput("Lists", "select a list to create task in", true),
	}
}

func (o *GetListOperation) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	accessToken := ctx.Auth.AccessToken
	input, err := sdk.InputToTypeSafely[getListProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	list, err := shared.GetList(accessToken, input.ListID)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (o *GetListOperation) Auth() *sdk.Auth {
	return nil
}

func (o *GetListOperation) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":      "list123",
		"name":    "Example List",
		"content": "List description",
		"statuses": []map[string]any{
			{
				"id":     "st123",
				"status": "Open",
				"color":  "#d3d3d3",
			},
		},
		"task_count": "24",
	}
}

func (o *GetListOperation) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetListOperation() sdk.Action {
	return &GetListOperation{}
}
