package triggers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type taskCreatedTriggerProps struct {
	ListID      string  `json:"list-id"`
	WorkspaceID *string `json:"workspace-id"`
}

type TaskCreatedTrigger struct{}

func (t *TaskCreatedTrigger) Name() string {
	return "Task Created"
}

func (t *TaskCreatedTrigger) Description() string {
	return "Triggers a workflow when a new task is created in a specified ClickUp workspace or list, allowing you to automate subsequent actions based on new task creation events."
}

func (t *TaskCreatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *TaskCreatedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &taskCreatedDocs,
	}
}

func (t *TaskCreatedTrigger) Icon() *string {
	icon := "material-symbols:add-task-outline"
	return &icon
}

func (t *TaskCreatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"list-id":      shared.GetListsInput("Lists", "select a list to get tasks from", true),
		"workspace-id": shared.GetWorkSpaceInput("Workspaces", "select a workspace", true),
	}
}

func (t *TaskCreatedTrigger) Start(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *TaskCreatedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *TaskCreatedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[taskCreatedTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if input.ListID == "" && (input.WorkspaceID == nil || *input.WorkspaceID == "") {
		return nil, errors.New("either List ID or Workspace ID must be provided")
	}

	lastRunTime := ctx.Metadata().LastRun

	var createdSince string
	if lastRunTime != nil {
		createdSince = strconv.FormatInt(lastRunTime.UnixNano()/int64(time.Millisecond), 10)
	}

	var endpoint string
	var queryParam string

	if input.ListID != "" {
		endpoint = "/v2/list/" + input.ListID + "/task"
	} else {
		endpoint = "/v2/team/" + *input.WorkspaceID + "/task"
	}

	if createdSince != "" {
		queryParam = fmt.Sprintf("?date_created_gt=%s", createdSince)
		endpoint += queryParam
	}

	response, err := shared.GetClickUpClient(ctx.Auth.AccessToken, endpoint, http.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	tasksArray, ok := response["tasks"].([]interface{})
	if !ok {
		return nil, errors.New("invalid response format: tasks field is not an array")
	}

	if len(tasksArray) == 0 {
		return []interface{}{}, nil
	}

	return tasksArray, nil
}

func (t *TaskCreatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *TaskCreatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *TaskCreatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":   "abc123",
		"name": "New Task",
		"status": map[string]any{
			"status": "Open",
			"color":  "#d3d3d3",
		},
		"date_created": "1647354847362",
		"creator": map[string]any{
			"id":       "123456",
			"username": "John Doe",
			"email":    "john@example.com",
		},
	}
}

func NewTaskCreatedTrigger() sdk.Trigger {
	return &TaskCreatedTrigger{}
}
