package triggers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/asana/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type taskUpdatedTriggerProps struct {
	ProjectID string `json:"project_id,omitempty"`
}

type TaskUpdatedTrigger struct{}

func (t *TaskUpdatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "task_updated",
		DisplayName:   "Task Updated",
		Description:   "Triggers a workflow whenever a task is updated in a specified Asana workspace or project.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: taskUpdatedDocs,
		Icon:          "simple-icons:asana",
		SampleOutput: map[string]any{
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
		},
	}
}

func (t *TaskUpdatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *TaskUpdatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("task_updated", "Task Updated")

	shared.RegisterProjectsProps(form)

	schema := form.Build()

	return schema
}

func (t *TaskUpdatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *TaskUpdatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *TaskUpdatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[taskUpdatedTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the last run time from metadata
	lastRun, err := ctx.GetMetadata("lastRun")
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	var updatedSince string
	if lastRunTime, ok := lastRun.(*time.Time); ok && lastRunTime != nil {
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

	response, err := shared.GetAsanaClient(authCtx.Token.AccessToken, url, http.MethodGet, nil)
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

func (t *TaskUpdatedTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewTaskUpdatedTrigger() sdk.Trigger {
	return &TaskUpdatedTrigger{}
}
