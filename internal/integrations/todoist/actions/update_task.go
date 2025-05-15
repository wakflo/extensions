package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/juicycleff/smartform/v1"
	fastshot "github.com/opus-domini/fast-shot"
	"github.com/wakflo/extensions/internal/integrations/todoist/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type updateTaskActionProps struct {
	shared.UpdateTask
	TaskID *string `json:"taskId"`
}

type UpdateTaskAction struct{}

func (a *UpdateTaskAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "update_task",
		DisplayName:   "Update Task",
		Description:   "Updates the status and details of an existing task in your workflow, allowing you to reflect changes or new information without having to recreate the task.",
		Type:          core.ActionTypeAction,
		Documentation: updateTaskDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *UpdateTaskAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("update_task", "Update Task")

	shared.RegisterTasksProps(form, "taskId", "Task ID", "ID of the task to update", true)

	form.TextareaField("content", "Content").
		Placeholder("The task's content. It may contain some markdown-formatted text and hyperlinks").
		HelpText("The task's content. It may contain some markdown-formatted text and hyperlinks")

	form.TextareaField("description", "Description").
		Placeholder("A description for the task. This value may contain some markdown-formatted text and hyperlinks.").
		HelpText("A description for the task. This value may contain some markdown-formatted text and hyperlinks.")

	labelsArray := form.ArrayField("labels", "Labels")
	labelGroup := labelsArray.ObjectTemplate("label", "")
	labelGroup.TextField("value", "Label").
		Placeholder("Label").
		Required(true).
		HelpText("Label")

	form.NumberField("priority", "Priority").
		Placeholder("Task priority from 1 (normal) to 4 (urgent).").
		HelpText("Task priority from 1 (normal) to 4 (urgent).")

	form.DateTimeField("dueDate", "Due date").
		Placeholder("Specific date in YYYY-MM-DD format relative to user's timezone").
		HelpText("Specific date in YYYY-MM-DD format relative to user's timezone")

	form.NumberField("duration", "Duration").
		Placeholder("A positive (greater than zero) integer for the amount of duration_unit the task will take, or null to unset. If specified, you must define a duration_unit.").
		HelpText("A positive (greater than zero) integer for the amount of duration_unit the task will take, or null to unset. If specified, you must define a duration_unit.")

	schema := form.Build()
	return schema
}

func (a *UpdateTaskAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateTaskActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	qu := fastshot.NewClient(shared.BaseAPI).
		Auth().BearerToken(authCtx.AccessToken).
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

func (a *UpdateTaskAction) Auth() *core.AuthMetadata {
	return nil
}

func NewUpdateTaskAction() sdk.Action {
	return &UpdateTaskAction{}
}
