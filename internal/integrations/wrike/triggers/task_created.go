package triggers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/wakflo/extensions/internal/integrations/wrike/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type taskCreatedTriggerProps struct {
	FolderID string `json:"folderId"`
	Limit    int    `json:"limit"`
}

type TaskCreatedTrigger struct{}

func (t *TaskCreatedTrigger) Name() string {
	return "Task Created"
}

func (t *TaskCreatedTrigger) Description() string {
	return "Triggers when a new task is created in your Wrike account."
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
	icon := "mdi:clipboard-plus"
	return &icon
}

func (t *TaskCreatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"folderId": shared.GetFoldersInput(),
		"limit": autoform.NewNumberField().
			SetDisplayName("Limit").
			SetDescription("Maximum number of tasks to return when triggered (1-100).").
			SetRequired(false).
			Build(),
	}
}

// Start initializes the trigger, required for event and webhook triggers in a lifecycle context.
func (t *TaskCreatedTrigger) Start(ctx sdk.LifecycleContext) error {
	// Not needed for polling triggers
	return nil
}

// Stop shuts down the trigger, cleaning up resources and performing necessary teardown operations.
func (t *TaskCreatedTrigger) Stop(ctx sdk.LifecycleContext) error {
	// Not needed for polling triggers
	return nil
}

// Execute performs the main logic of the trigger by checking for new tasks since the last run.
func (t *TaskCreatedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[taskCreatedTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	lastRunTime := ctx.Metadata().LastRun

	var createdTime string
	if lastRunTime != nil {
		createdTime = lastRunTime.UTC().Format(time.RFC3339)
	} else {
		createdTime = ""
	}

	var endpoint string
	if input.FolderID != "" {
		endpoint = fmt.Sprintf("/folders/%s/tasks", input.FolderID)
	} else {
		endpoint = "/tasks"
	}

	limit := 25
	if input.Limit > 0 {
		limit = input.Limit
	}

	endpoint = fmt.Sprintf("%s?createdDate={\"start\":\"%s\"}&sortField=CreatedDate&sortOrder=Asc&limit=%d",
		endpoint, createdTime, limit)

	response, err := shared.GetWrikeClient(ctx.Auth.AccessToken, endpoint)
	if err != nil {
		return nil, errors.New("invalid response format from Wrike API")
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
	return []interface{}{
		map[string]interface{}{
			"id":          "IEADTSKYA5CKABNW",
			"accountId":   "IEADTSKY",
			"title":       "New Task",
			"description": "This is a new task created in Wrike",
			"status":      "Active",
			"importance":  "Normal",
			"createdDate": "2023-03-20T14:30:45.000Z",
			"updatedDate": "2023-03-20T14:30:45.000Z",
		},
	}
}

func NewTaskCreatedTrigger() sdk.Trigger {
	return &TaskCreatedTrigger{}
}
