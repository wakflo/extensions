package triggers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/wrike/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type taskCreatedTriggerProps struct {
	FolderID string `json:"folderId"`
	Limit    int    `json:"limit"`
}

type TaskCreatedTrigger struct{}

func (t *TaskCreatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "task_created",
		DisplayName:   "Task Created",
		Description:   "Triggers when a new task is created in your Wrike account.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: taskCreatedDocs,
		SampleOutput: []interface{}{
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
		},
	}
}

func (t *TaskCreatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("task_created", "Task Created")

	shared.GetFoldersProp(form)
	form.NumberField("limit", "Limit").
		Placeholder("Enter limit").
		Required(false).
		HelpText("Maximum number of tasks to return when triggered (1-100).")

	schema := form.Build()

	return schema
}

// Start initializes the trigger, required for event and webhook triggers in a lifecycle context.
func (t *TaskCreatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	// Not needed for polling triggers
	return nil
}

// Stop shuts down the trigger, cleaning up resources and performing necessary teardown operations.
func (t *TaskCreatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	// Not needed for polling triggers
	return nil
}

// Execute performs the main logic of the trigger by checking for new tasks since the last run.
func (t *TaskCreatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[taskCreatedTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	lastRunTime, err := ctx.GetMetadata("lastrun")

	var createdTime string
	if lastRunTime != nil {
		createdTime = lastRunTime.(*time.Time).UTC().Format(time.RFC3339)
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

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	response, err := shared.GetWrikeClient(token, endpoint)
	if err != nil {
		return nil, errors.New("invalid response format from Wrike API")
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
