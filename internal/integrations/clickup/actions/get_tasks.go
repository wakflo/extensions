package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

// GetTasksOperation structure and methods
type getTasksProps struct {
	WorkspaceID string `json:"workspace-id"`
	SpaceID     string `json:"space-id"`
	FolderID    string `json:"folder-id"`
	ListID      string `json:"list-id"`
}

type GetTasksOperation struct{}

// Metadata returns metadata about the action
func (o *GetTasksOperation) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_tasks",
		DisplayName:   "Get Tasks",
		Description:   "Retrieves all tasks from a specified ClickUp list.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: getTasksDocs,
		Icon:          "material-symbols:task-alt-outline",
		SampleOutput: map[string]any{
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
		},
		Settings: sdkcore.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (o *GetTasksOperation) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_tasks", "Get Tasks")

	shared.RegisterWorkSpaceInput(form, "Workspace", "select a workspace", true)
	shared.RegisterSpacesInput(form, "Space", "select a space", true)
	shared.RegisterFoldersInput(form, "Folder", "select a folder", true)
	shared.RegisterListsInput(form, "List", "select a list", true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (o *GetTasksOperation) Auth() *sdkcore.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (o *GetTasksOperation) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getTasksProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	accessToken := authCtx.Token.AccessToken
	url := "/v2/list/" + input.ListID + "/task"

	tasks, err := shared.GetData(accessToken, url)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func NewGetTasksAction() sdk.Action {
	return &GetTasksOperation{}
}
