package actions

import (
	"encoding/json"
	"errors"
	"io"

	fastshot "github.com/opus-domini/fast-shot"
	"github.com/wakflo/extensions/internal/integrations/todoist/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getActiveTaskActionProps struct {
	TaskID string `json:"taskId"`
}

type GetActiveTaskAction struct{}

func (a *GetActiveTaskAction) Name() string {
	return "Get Active Task"
}

func (a *GetActiveTaskAction) Description() string {
	return "Retrieves the currently active task in the workflow, allowing you to access and manipulate its properties or trigger subsequent actions based on its status."
}

func (a *GetActiveTaskAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetActiveTaskAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getActiveTaskDocs,
	}
}

func (a *GetActiveTaskAction) Icon() *string {
	return nil
}

func (a *GetActiveTaskAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"taskId": shared.GetTasksInput("Task ID", "ID of the active task you want to retrieve", true),
	}
}

func (a *GetActiveTaskAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getActiveTaskActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	qu := fastshot.NewClient(shared.BaseAPI).
		Auth().BearerToken(ctx.Auth.AccessToken).
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

func (a *GetActiveTaskAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetActiveTaskAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *GetActiveTaskAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetActiveTaskAction() sdk.Action {
	return &GetActiveTaskAction{}
}
