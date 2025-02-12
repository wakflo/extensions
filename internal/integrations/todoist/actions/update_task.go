package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	fastshot "github.com/opus-domini/fast-shot"
	"github.com/wakflo/extensions/internal/integrations/todoist/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type updateTaskActionProps struct {
	shared.UpdateTask
	TaskID *string `json:"taskId"`
}

type UpdateTaskAction struct{}

func (a *UpdateTaskAction) Name() string {
	return "Update Task"
}

func (a *UpdateTaskAction) Description() string {
	return "Updates the status and details of an existing task in your workflow, allowing you to reflect changes or new information without having to recreate the task."
}

func (a *UpdateTaskAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *UpdateTaskAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &updateTaskDocs,
	}
}

func (a *UpdateTaskAction) Icon() *string {
	return nil
}

func (a *UpdateTaskAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"taskId": shared.GetTasksInput("Task ID", "ID of the task to update", true),

		"content": autoform.NewMarkdownField().
			SetDisplayName("Content").
			SetDescription("The task's content. It may contain some markdown-formatted text and hyperlinks").
			SetRequired(false).Build(),

		"description": autoform.NewMarkdownField().
			SetDisplayName("Description").
			SetDescription("A description for the task. This value may contain some markdown-formatted text and hyperlinks.").
			SetRequired(false).Build(),

		"labels": autoform.NewArrayField().
			SetDisplayName("Labels").
			SetDescription("The task's labels (a list of names that may represent either personal or shared labels)").
			SetItems(
				autoform.NewShortTextField().
					SetDisplayName("Label").
					SetDescription("Label").
					SetRequired(true).
					Build(),
			).
			SetRequired(false).Build(),

		"priority": autoform.NewNumberField().
			SetDisplayName("Priority").
			SetDescription("Task priority from 1 (normal) to 4 (urgent).").
			SetRequired(false).Build(),

		"dueDate": autoform.NewDateTimeField().
			SetDisplayName("Due date").
			SetDescription("Specific date in YYYY-MM-DD format relative to user's timezone").
			SetRequired(false).Build(),

		"duration": autoform.NewNumberField().
			SetDisplayName("Duration").
			SetDescription("A positive (greater than zero) integer for the amount of duration_unit the task will take, or null to unset. If specified, you must define a duration_unit.").
			SetRequired(false).Build(),
	}
}

func (a *UpdateTaskAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateTaskActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	qu := fastshot.NewClient(shared.BaseAPI).
		Auth().BearerToken(ctx.Auth.AccessToken).
		Header().
		AddAccept("application/json").
		Build().POST(fmt.Sprintf("/tasks/%v", input.TaskID)).Body().AsJSON(input.UpdateTask)

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

func (a *UpdateTaskAction) Auth() *sdk.Auth {
	return nil
}

func (a *UpdateTaskAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *UpdateTaskAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewUpdateTaskAction() sdk.Action {
	return &UpdateTaskAction{}
}
