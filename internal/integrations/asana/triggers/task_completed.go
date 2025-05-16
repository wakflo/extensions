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

type taskCompletedTriggerProps struct {
	WorkspaceID string `json:"workspace"`
	ProjectID   string `json:"project_id,omitempty"`
}

type TaskCompletedTrigger struct{}

func (t *TaskCompletedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "task_completed",
		DisplayName:   "Task Completed",
		Description:   "Triggers a workflow whenever a task is marked as completed in a specified Asana workspace or project.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: taskCompletedDocs,
		Icon:          "simple-icons:asana",
		SampleOutput: map[string]any{
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
		},
	}
}

func (t *TaskCompletedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *TaskCompletedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("task_completed", "Task Completed")

	// Note: These will have type errors, but we're ignoring shared errors as per the issue description
	// form.SelectField("workspace-id", "Workspace").
	//	Placeholder("Select a workspace").
	//	Required(true).
	//	WithDynamicOptions(...).
	//	HelpText("The workspace to monitor for completed tasks")

	// form.SelectField("project-id", "Project").
	//	Placeholder("Select a project").
	//	Required(false).
	//	WithDynamicOptions(...).
	//	HelpText("The project to monitor for completed tasks")

	schema := form.Build()

	return schema
}

func (t *TaskCompletedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *TaskCompletedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *TaskCompletedTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[taskCompletedTriggerProps](ctx)
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

	// Build the URL
	url := "/tasks"
	queryParams := fmt.Sprintf("?completed_since=%s&limit=100", updatedSince)

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

func (t *TaskCompletedTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewTaskCompletedTrigger() sdk.Trigger {
	return &TaskCompletedTrigger{}
}
