package actions

import (
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

// GetTasksOperation structure and methods
type getTasksProps struct {
	WorkspaceID string `json:"workspace-id"`
	SpaceID     string `json:"space-id"`
	FolderID    string `json:"folder-id"`
	ListID      string `json:"list-id"`
}

type GetTasksOperation struct{}

func (o *GetTasksOperation) Name() string {
	return "Get Tasks"
}

func (o *GetTasksOperation) Description() string {
	return "Retrieves all tasks from a specified ClickUp list."
}

func (o *GetTasksOperation) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (o *GetTasksOperation) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getTasksDocs,
	}
}

func (o *GetTasksOperation) Icon() *string {
	icon := "material-symbols:task-alt-outline"
	return &icon
}

func (o *GetTasksOperation) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"workspace-id": shared.GetWorkSpaceInput("Workspaces", "select a workspace", true),
		"space-id":     shared.GetSpacesInput("Spaces", "select a space", true),
		"folder-id":    shared.GetFoldersInput("Folders", "select a folder", true),
		"list-id":      shared.GetListsInput("Lists", "select a list to get tasks from", true),
	}
}

func (o *GetTasksOperation) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	accessToken := ctx.Auth.AccessToken
	input, err := sdk.InputToTypeSafely[getTasksProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}
	url := "/v2/list/" + input.ListID + "/task"

	tasks, err := shared.GetData(accessToken, url)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (o *GetTasksOperation) Auth() *sdk.Auth {
	return nil
}

func (o *GetTasksOperation) SampleData() sdkcore.JSON {
	return map[string]any{
		"tasks": []map[string]any{
			{
				"id":   "abc123",
				"name": "Example Task 1",
				"status": map[string]string{
					"status": "Open",
					"color":  "#d3d3d3",
				},
			},
			{
				"id":   "def456",
				"name": "Example Task 2",
				"status": map[string]string{
					"status": "In Progress",
					"color":  "#4286f4",
				},
			},
		},
	}
}

func (o *GetTasksOperation) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetTasksOperation() sdk.Action {
	return &GetTasksOperation{}
}
