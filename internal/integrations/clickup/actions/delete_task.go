package actions

import (
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type deleteTaskProps struct {
	TaskID string `json:"task-id"`
}

type DeleteTaskOperation struct{}

// Metadata returns metadata about the action
func (o *DeleteTaskOperation) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "delete_task",
		DisplayName:   "Delete Task",
		Description:   "Deletes a task from ClickUp.",
		Type:          core.ActionTypeAction,
		Documentation: deleteTaskDocs,
		Icon:          "material-symbols:delete",
		SampleOutput: map[string]any{
			"success": true,
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (o *DeleteTaskOperation) Properties() *smartform.FormSchema {
	form := smartform.NewForm("delete_task", "Delete Task")

	form.TextField("task-id", "task-id").
		Placeholder("Task ID").
		HelpText("The ID of the task to delete").
		Required(true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (o *DeleteTaskOperation) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (o *DeleteTaskOperation) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[deleteTaskProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	accessToken := authCtx.Token.AccessToken
	reqURL := shared.BaseURL + "/v2/task/" + input.TaskID
	req, err := http.NewRequest(http.MethodDelete, reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", accessToken)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return map[string]interface{}{
		"task": "Task Deleted",
	}, nil
}

func NewDeleteTaskOperation() sdk.Action {
	return &DeleteTaskOperation{}
}
