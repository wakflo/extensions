package triggers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type taskUpdatedTriggerProps struct {
	ListID      string  `json:"list-id"`
	WorkspaceID *string `json:"workspace-id"`
	Status      *string `json:"status"`
}

type TaskUpdatedTrigger struct{}

func (t *TaskUpdatedTrigger) Name() string {
	return "Task Updated"
}

func (t *TaskUpdatedTrigger) Description() string {
	return "Triggers a workflow when an existing task is updated in ClickUp, including changes to status, priority, assignees, or due dates, enabling automated reactions to task modification events."
}

func (t *TaskUpdatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *TaskUpdatedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &taskUpdatedDocs,
	}
}

func (t *TaskUpdatedTrigger) Icon() *string {
	icon := "material-symbols:update"
	return &icon
}

func (t *TaskUpdatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"workspace-id": shared.GetWorkSpaceInput("Workspaces", "select a workspace", true),
		"list-id":      shared.GetListsInput("Lists", "select a list to get tasks from", true),
		"status": autoform.NewShortTextField().
			SetDisplayName("Status Filter").
			SetDescription("Only trigger when tasks are updated to this status. Leave empty to trigger on any status change.").
			SetRequired(false).Build(),
	}
}

func (t *TaskUpdatedTrigger) Start(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *TaskUpdatedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *TaskUpdatedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[taskUpdatedTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if input.ListID == "" && (input.WorkspaceID == nil || *input.WorkspaceID == "") {
		return nil, errors.New("either List ID or Workspace ID must be provided")
	}

	lastRunTime := ctx.Metadata().LastRun

	var updatedSince string
	if lastRunTime != nil {
		updatedSince = strconv.FormatInt(lastRunTime.UnixNano()/int64(time.Millisecond), 10)
	}

	var endpoint string
	var queryParams string

	if input.ListID != "" {
		endpoint = fmt.Sprintf("/list/%s/task", input.ListID)
	} else {
		endpoint = fmt.Sprintf("/team/%s/task", *input.WorkspaceID)
	}

	if updatedSince != "" {
		queryParams = "?date_updated_gt=" + updatedSince

		if input.Status != nil && *input.Status != "" {
			queryParams += "&statuses=" + *input.Status
		}

		endpoint += queryParams
	}

	response, err := shared.GetClickUpClient(ctx.Auth.AccessToken, endpoint, http.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	tasksArray, ok := response["tasks"].([]interface{})
	if !ok {
		return nil, errors.New("invalid response format: tasks field is not an array")
	}

	var updatedTasks []interface{}
	if lastRunTime != nil {
		for _, task := range tasksArray {
			taskMap, ok := task.(map[string]interface{})
			if !ok {
				continue
			}

			dateCreatedStr, createdOk := taskMap["date_created"].(string)
			dateUpdatedStr, updatedOk := taskMap["date_updated"].(string)

			if !createdOk || !updatedOk {
				continue
			}

			var dateCreated, dateUpdated int64
			fmt.Sscanf(dateCreatedStr, "%d", &dateCreated)
			fmt.Sscanf(dateUpdatedStr, "%d", &dateUpdated)

			if dateUpdated > dateCreated {
				updatedTasks = append(updatedTasks, task)
			}
		}
	} else {
		updatedTasks = tasksArray
	}

	if len(updatedTasks) == 0 {
		return []interface{}{}, nil
	}

	return updatedTasks, nil
}

func (t *TaskUpdatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *TaskUpdatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *TaskUpdatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":   "abc123",
		"name": "Updated Task",
		"status": map[string]any{
			"status": "In Progress",
			"color":  "#4194f6",
		},
		"date_created": "1647354847362",
		"date_updated": "1647354987362",
		"update_fields": []string{
			"status",
			"assignees",
			"due_date",
		},
		"assignees": []map[string]any{
			{
				"id":       "123456",
				"username": "John Doe",
				"email":    "john@example.com",
			},
		},
	}
}

func NewTaskUpdatedTrigger() sdk.Trigger {
	return &TaskUpdatedTrigger{}
}
