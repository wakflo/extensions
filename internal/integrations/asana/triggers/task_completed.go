package triggers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/wakflo/extensions/internal/integrations/asana/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type taskCompletedTriggerProps struct {
	WorkspaceID string `json:"workspace-id"`
	ProjectID   string `json:"project-id,omitempty"`
}

type TaskCompletedTrigger struct{}

func (t *TaskCompletedTrigger) Name() string {
	return "Task Completed"
}

func (t *TaskCompletedTrigger) Description() string {
	return "Triggers a workflow whenever a task is marked as completed in a specified Asana workspace or project."
}

func (t *TaskCompletedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *TaskCompletedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &taskCompletedDocs,
	}
}

func (t *TaskCompletedTrigger) Icon() *string {
	icon := "simple-icons:asana"
	return &icon
}

func (t *TaskCompletedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"workspace-id": shared.GetWorkspacesInput(),
		"project-id":   shared.GetProjectsInput(),
	}
}

func (t *TaskCompletedTrigger) Start(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *TaskCompletedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *TaskCompletedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[taskCompletedTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	lastRunTime := ctx.Metadata().LastRun

	var updatedSince string
	if lastRunTime != nil {
		updatedSince = lastRunTime.UTC().Format(time.RFC3339)
	} else {
		updatedSince = ""
	}

	// Build the URL
	url := "/tasks"
	queryParams := fmt.Sprintf("?completed_since=%s&limit=100", updatedSince)

	if input.ProjectID != "" {
		queryParams += fmt.Sprintf("&project=%s", input.ProjectID)
	}

	url += queryParams

	response, err := shared.GetAsanaClient(ctx.Auth.AccessToken, url, http.MethodGet, nil)
	if err != nil {
		return nil, errors.New("error fetching data")
	}

	data, ok := response["data"].([]interface{})
	if !ok {
		return nil, errors.New("invalid response format: data field is not an array")
	}

	var completedTasks []interface{}
	for _, task := range data {
		taskMap, ok := task.(map[string]interface{})
		if !ok {
			continue
		}

		completed, ok := taskMap["completed"].(bool)
		if !ok {
			continue
		}

		if completed {
			completedTasks = append(completedTasks, task)
		}
	}

	// If no completed tasks, return empty array
	if len(completedTasks) == 0 {
		return map[string]interface{}{
			"data": []interface{}{},
		}, nil
	}

	return map[string]interface{}{
		"data": completedTasks,
	}, nil
}

func (t *TaskCompletedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *TaskCompletedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *TaskCompletedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"data": []map[string]any{
			{
				"gid":           "1234567890",
				"name":          "Completed Task",
				"resource_type": "task",
				"created_at":    "2023-01-10T08:00:00.000Z",
				"completed":     true,
				"completed_at":  "2023-01-15T09:30:00.000Z",
				"notes":         "This task has been completed",
			},
		},
	}
}

func NewTaskCompletedTrigger() sdk.Trigger {
	return &TaskCompletedTrigger{}
}
