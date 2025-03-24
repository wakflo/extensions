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

type taskUpdatedTriggerProps struct {
	ProjectID string `json:"project-id,omitempty"`
}

type TaskUpdatedTrigger struct{}

func (t *TaskUpdatedTrigger) Name() string {
	return "Task Updated"
}

func (t *TaskUpdatedTrigger) Description() string {
	return "Triggers a workflow whenever a task is updated in a specified Asana workspace or project."
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
	icon := "simple-icons:asana"
	return &icon
}

func (t *TaskUpdatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"project-id": shared.GetProjectsInput(),
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

	lastRunTime := ctx.Metadata().LastRun

	var updatedSince string
	if lastRunTime != nil {
		updatedSince = lastRunTime.UTC().Format(time.RFC3339)
	} else {
		updatedSince = ""
	}

	url := "/tasks"
	queryParams := fmt.Sprintf("?modified_since=%s&limit=100", updatedSince)

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

	var updatedTasks []interface{}
	for _, task := range data {
		taskMap, ok := task.(map[string]interface{})
		if !ok {
			continue
		}

		createdAt, ok := taskMap["created_at"].(string)
		if !ok {
			continue
		}

		modifiedAt, ok := taskMap["modified_at"].(string)
		if !ok {
			continue
		}

		createdTime, _ := time.Parse(time.RFC3339, createdAt)
		modifiedTime, _ := time.Parse(time.RFC3339, modifiedAt)

		if modifiedTime.Sub(createdTime).Seconds() > 5 {
			updatedTasks = append(updatedTasks, task)
		}
	}

	if len(updatedTasks) == 0 {
		return map[string]interface{}{
			"data": []interface{}{},
		}, nil
	}

	return map[string]interface{}{
		"data": updatedTasks,
	}, nil
}

func (t *TaskUpdatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *TaskUpdatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *TaskUpdatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"data": []map[string]any{
			{
				"gid":           "1234567890",
				"name":          "Updated Task",
				"resource_type": "task",
				"created_at":    "2023-01-10T08:00:00.000Z",
				"modified_at":   "2023-01-15T09:30:00.000Z",
				"notes":         "This task has been updated with new information",
				"completed":     false,
			},
		},
	}
}

func NewTaskUpdatedTrigger() sdk.Trigger {
	return &TaskUpdatedTrigger{}
}
