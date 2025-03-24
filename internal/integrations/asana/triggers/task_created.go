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

type taskCreatedTriggerProps struct {
	WorkspaceID string `json:"workspace-id"`
	ProjectID   string `json:"project-id,omitempty"`
}

type TaskCreatedTrigger struct{}

func (t *TaskCreatedTrigger) Name() string {
	return "Task Created"
}

func (t *TaskCreatedTrigger) Description() string {
	return "Triggers a workflow whenever a new task is created in a specified Asana workspace or project."
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
	icon := "simple-icons:asana"
	return &icon
}

func (t *TaskCreatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"project-id": shared.GetProjectsInput(),
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

	lastRunTime := ctx.Metadata().LastRun

	var updatedSince string
	if lastRunTime != nil {
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

	response, err := shared.GetAsanaClient(ctx.Auth.AccessToken, url, http.MethodGet, nil)
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

func (t *TaskCreatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *TaskCreatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
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
	}
}

func NewTaskCreatedTrigger() sdk.Trigger {
	return &TaskCreatedTrigger{}
}
