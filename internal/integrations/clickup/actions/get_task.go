package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getTaskActionProps struct {
	TaskID string `json:"task-id"`
}

type GetTaskAction struct{}

// Metadata returns metadata about the action
func (a *GetTaskAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_task",
		DisplayName:   "Get Task",
		Description:   "Retrieves details of a specific ClickUp task by ID.",
		Type:          core.ActionTypeAction,
		Documentation: getTaskDocs,
		SampleOutput: map[string]any{
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
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GetTaskAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_task", "Get Task")

	shared.RegisterWorkSpaceInput(form, "Workspaces", "select a workspace", true)
	shared.RegisterSpacesInput(form, "Spaces", "select a space", true)
	shared.RegisterFoldersInput(form, "Folders", "select a folder", true)
	shared.RegisterListsInput(form, "Lists", "select a list to create task in", true)
	shared.RegisterTasksInput(form, "Tasks", "select a task to update", true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetTaskAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GetTaskAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[getTaskActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the token source from the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	url := "/v2/task/" + input.TaskID

	tasks, err := shared.GetData(authCtx.AccessToken, url)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func NewGetTaskAction() sdk.Action {
	return &GetTaskAction{}
}
