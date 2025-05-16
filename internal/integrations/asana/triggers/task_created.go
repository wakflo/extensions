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

type taskCreatedTriggerProps struct {
	WorkspaceID string `json:"workspace_id"`
	ProjectID   string `json:"project_id,omitempty"`
}

type TaskCreatedTrigger struct{}

func (t *TaskCreatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "task_created",
		DisplayName:   "Task Created",
		Description:   "Triggers a workflow whenever a new task is created in a specified Asana workspace or project.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: taskCreatedDocs,
		Icon:          "simple-icons:asana",
		SampleOutput: map[string]any{
			"data": []map[string]any{
				{
					"gid":           "1234567890",
					"name":          "New Task",
					"resource_type": "task",
					"created_at":    "2023-01-15T08:00:00.000Z",
					"notes":         "This is a new task that was just created",
					"completed":     false,
				},
			},
		},
	}
}

func (t *TaskCreatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *TaskCreatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("task_created", "Task Created")

	shared.RegisterProjectsProps(form)

	schema := form.Build()

	return schema
}

func (t *TaskCreatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *TaskCreatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *TaskCreatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[taskCreatedTriggerProps](ctx)
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
	queryParams := fmt.Sprintf("?created_since=%s&limit=100", updatedSince)

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

	if len(data) == 0 {
		return map[string]interface{}{
			"data": []interface{}{},
		}, nil
	}

	return response, nil
}

func (t *TaskCreatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *TaskCreatedTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewTaskCreatedTrigger() sdk.Trigger {
	return &TaskCreatedTrigger{}
}
