package actions

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/juicycleff/smartform/v1"
	fastshot "github.com/opus-domini/fast-shot"
	"github.com/wakflo/extensions/internal/integrations/todoist/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createTaskActionProps struct {
	Content     string     `json:"content"`
	Description *string    `json:"description"`
	ProjectID   *string    `json:"project_id"`
	SectionID   *string    `json:"section_id"`
	ParentID    *string    `json:"parent_id"`
	Labels      []string   `json:"labels"`
	Order       *string    `json:"order"`
	Priority    *int       `json:"priority"`
	DueDate     *time.Time `json:"dueDate"`
}

type CreateTaskAction struct{}

func (a *CreateTaskAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_task",
		DisplayName:   "Create Task",
		Description:   "Create Task: Automatically generates and assigns a new task to a team member or group, allowing you to streamline workflows and ensure timely completion of tasks.",
		Type:          core.ActionTypeAction,
		Documentation: createTaskDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *CreateTaskAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_task", "Create Task")
	shared.RegisterProjectsProps(form)
	shared.RegisterSectionsProps(form)
	shared.RegisterTasksProps(form, "parent_id", "Parent Task ID", "Parent task ID.", false)
	form.TextareaField("content", "Content").
		Placeholder("The task's content. It may contain some markdown-formatted text and hyperlinks").
		Required(true).
		HelpText("The task's content. It may contain some markdown-formatted text and hyperlinks")

	form.TextareaField("description", "Description").
		Placeholder("A description for the task. This value may contain some markdown-formatted text and hyperlinks.").
		HelpText("A description for the task. This value may contain some markdown-formatted text and hyperlinks.")

	form.NumberField("order", "Order").
		Placeholder("Non-zero integer value used by clients to sort tasks under the same parent.").
		HelpText("Non-zero integer value used by clients to sort tasks under the same parent.")

	labelsArray := form.ArrayField("labels", "Labels")
	labelGroup := labelsArray.ObjectTemplate("label", "")
	labelGroup.TextField("value", "Label").
		Placeholder("Label").
		Required(true).
		HelpText("The task's labels (a list of names that may represent either personal or shared labels)")

	form.NumberField("priority", "Priority").
		Placeholder("Task priority from 1 (normal) to 4 (urgent).").
		HelpText("Task priority from 1 (normal) to 4 (urgent).")

	form.DateTimeField("dueDate", "Due date").
		Placeholder("Specific date in YYYY-MM-DD format relative to user's timezone").
		HelpText("Specific date in YYYY-MM-DD format relative to user's timezone")

	schema := form.Build()
	return schema
}

func (a *CreateTaskAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createTaskActionProps](ctx)
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
		Build().POST("/tasks").Body().AsJSON(input)

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

func (a *CreateTaskAction) Auth() *core.AuthMetadata {
	return nil
}

func NewCreateTaskAction() sdk.Action {
	return &CreateTaskAction{}
}
