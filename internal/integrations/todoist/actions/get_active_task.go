package actions

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/juicycleff/smartform/v1"
	fastshot "github.com/opus-domini/fast-shot"
	"github.com/wakflo/extensions/internal/integrations/todoist/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getActiveTaskActionProps struct {
	TaskID string `json:"taskId"`
}

type GetActiveTaskAction struct{}

func (a *GetActiveTaskAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_active_task",
		DisplayName:   "Get Active Task",
		Description:   "Retrieves the currently active task in the workflow, allowing you to access and manipulate its properties or trigger subsequent actions based on its status.",
		Type:          core.ActionTypeAction,
		Documentation: getActiveTaskDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *GetActiveTaskAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_active_task", "Get Active Task")

	shared.RegisterTasksProps(form, "taskId", "Task ID", "ID of the active task you want to retrieve", true)

	schema := form.Build()
	return schema
}

func (a *GetActiveTaskAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getActiveTaskActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	qu := fastshot.NewClient(shared.BaseAPI).
		Auth().BearerToken(authCtx.Token.AccessToken).
		Header().
		AddAccept("application/json").
		Build().GET("/tasks/" + input.TaskID)

	rsp, err := qu.Send()
	if err != nil {
		return nil, err
	}

	if rsp.Status().IsError() {
		return nil, errors.New(rsp.Status().Text())
	}

	bytes, err := io.ReadAll(rsp.Raw().Body)
	if err != nil {
		return nil, err
	}

	var task shared.Task
	err = json.Unmarshal(bytes, &task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (a *GetActiveTaskAction) Auth() *core.AuthMetadata {
	return nil
}

func NewGetActiveTaskAction() sdk.Action {
	return &GetActiveTaskAction{}
}
