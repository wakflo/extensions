package actions

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	fastshot "github.com/opus-domini/fast-shot"
	"github.com/wakflo/extensions/internal/integrations/todoist/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
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

func (a *CreateTaskAction) Name() string {
	return "Create Task"
}

func (a *CreateTaskAction) Description() string {
	return "Create Task: Automatically generates and assigns a new task to a team member or group, allowing you to streamline workflows and ensure timely completion of tasks."
}

func (a *CreateTaskAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateTaskAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createTaskDocs,
	}
}

func (a *CreateTaskAction) Icon() *string {
	return nil
}

func (a *CreateTaskAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"content": autoform.NewMarkdownField().
			SetDisplayName("Content").
			SetDescription("The task's content. It may contain some markdown-formatted text and hyperlinks").
			SetRequired(true).Build(),

		"description": autoform.NewMarkdownField().
			SetDisplayName("Description").
			SetDescription("A description for the task. This value may contain some markdown-formatted text and hyperlinks.").
			SetRequired(false).Build(),

		"project_id": shared.GetProjectsInput(),
		"section_id": shared.GetSectionsInput(),
		"parent_id":  shared.GetTasksInput("Parent Task ID", "Parent task ID.", false),
		"order": autoform.NewNumberField().
			SetDisplayName("Order").
			SetDescription("Non-zero integer value used by clients to sort tasks under the same parent.").
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
	}
}

func (a *CreateTaskAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createTaskActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	qu := fastshot.NewClient(shared.BaseAPI).
		Auth().BearerToken(ctx.Auth.AccessToken).
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

func (a *CreateTaskAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateTaskAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateTaskAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateTaskAction() sdk.Action {
	return &CreateTaskAction{}
}
