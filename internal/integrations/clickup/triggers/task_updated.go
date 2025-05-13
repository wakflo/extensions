package triggers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type taskUpdatedTriggerProps struct {
	ListID      string  `json:"list-id"`
	WorkspaceID *string `json:"workspace-id"`
	Status      *string `json:"status"`
}

type TaskUpdatedTrigger struct{}

func (t *TaskUpdatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "task_updated",
		DisplayName:   "Task Updated",
		Description:   "Triggers a workflow when an existing task is updated in ClickUp, including changes to status, priority, assignees, or due dates, enabling automated reactions to task modification events.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: taskUpdatedDocs,
		Icon:          "material-symbols:update",
		SampleOutput: map[string]any{
			"id":   "abc123",
			"name": "Updated Task",
			"status": map[string]any{
				"status": "In Progress",
				"color":  "#4194f6",
			},
		},
	}
}

func (t *TaskUpdatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("get_task", "Get Task")

	shared.RegisterWorkSpaceInput(form, "Workspaces", "select a workspace", true)
	shared.RegisterListsInput(form, "Lists", "select a list to create task in", true)

	form.TextField("status", "Status Filter").
		Placeholder("Leave empty to trigger on any status change.").
		Required(false).
		HelpText("Only trigger when tasks are updated to this status. Leave empty to trigger on any status change.").
		Build()

	schema := form.Build()

	return schema
}

func (t *TaskUpdatedTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func (t *TaskUpdatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *TaskUpdatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *TaskUpdatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[taskUpdatedTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	if input.ListID == "" && (input.WorkspaceID == nil || *input.WorkspaceID == "") {
		return nil, errors.New("either List ID or Workspace ID must be provided")
	}

	lastRunTime := ctx.LastRun()

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

	response, err := shared.GetClickUpClient(ctx.Auth().AccessToken, endpoint, http.MethodGet, nil)
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

func NewTaskUpdatedTrigger() sdk.Trigger {
	return &TaskUpdatedTrigger{}
}
