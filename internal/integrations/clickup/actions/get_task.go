package actions

import (
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getTaskProps struct {
	TaskID string `json:"task-id"`
}

type GetTaskOperation struct{}

func (o *GetTaskOperation) Name() string {
	return "Get Task"
}

func (o *GetTaskOperation) Description() string {
	return "Retrieves details of a specific ClickUp task by ID."
}

func (o *GetTaskOperation) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (o *GetTaskOperation) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getTaskDocs,
	}
}

func (o *GetTaskOperation) Icon() *string {
	icon := "material-symbols:task"
	return &icon
}

func (o *GetTaskOperation) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"workspace-id": shared.GetWorkSpaceInput("Workspaces", "select a workspace", true),
		"space-id":     shared.GetSpacesInput("Spaces", "select a space", true),
		"folder-id":    shared.GetFoldersInput("Folders", "select a folder", true),
		"list-id":      shared.GetListsInput("Lists", "select a list to create task in", true),
		"task-id":      shared.GetTasksInput("Tasks", "select a task to update", true),
	}
}

func (o *GetTaskOperation) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	accessToken := ctx.Auth.AccessToken
	input, err := sdk.InputToTypeSafely[getTaskProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	url := "/v2/task/" + input.TaskID

	tasks, err := shared.GetData(accessToken, url)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (o *GetTaskOperation) Auth() *sdk.Auth {
	return nil
}

func (o *GetTaskOperation) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":          "abc123",
		"name":        "Example Task",
		"description": "This is a sample task",
		"status": map[string]string{
			"status": "Open",
			"color":  "#d3d3d3",
		},
		"priority": map[string]any{
			"priority": "High",
			"color":    "#f50000",
		},
		"date_created": "1647354847362",
		"date_updated": "1647354847362",
		"creator": map[string]any{
			"id":       "123456",
			"username": "John Doe",
			"email":    "john@example.com",
		},
	}
}

func (o *GetTaskOperation) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetTaskOperation() sdk.Action {
	return &GetTaskOperation{}
}
