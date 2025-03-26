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

type taskUpdatedTriggerProps struct {
	FolderID string `json:"folderId"`
	Limit    int    `json:"limit"`
}

type TaskUpdatedTrigger struct{}

func (t *TaskUpdatedTrigger) Name() string {
	return "Task Updated"
}

func (t *TaskUpdatedTrigger) Description() string {
	return "Triggers when a task is updated in your Wrike account."
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
	icon := "mdi:clipboard-edit"
	return &icon
}

func (t *TaskUpdatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"folderId": shared.GetFoldersInput(),
		"limit": autoform.NewNumberField().
			SetDisplayName("Limit").
			SetDescription("Maximum number of tasks to return when triggered (1-100).").
			SetRequired(false).
			Build(),
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

	var updatedTime string
	if lastRunTime != nil {
		updatedTime = lastRunTime.UTC().Format(time.RFC3339)
	} else {
		updatedTime = ""
	}

	// Construct the API endpoint
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

	endpoint = fmt.Sprintf("%s?updatedDate={\"start\":\"%s\"}&sortField=UpdatedDate&sortOrder=Asc&limit=%d",
		endpoint, updatedTime, limit)

	response, err := shared.GetWrikeClient(ctx.Auth.AccessToken, endpoint)
	if err != nil {
		return nil, errors.New("error fetching data from Wrike API")
	}

	responseMap, ok := response.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid response format from Wrike API")
	}

	data, ok := responseMap["data"].([]interface{})
	if !ok {
		return nil, errors.New("invalid response format from Wrike API")
	}

	if len(data) == 0 {
		return []interface{}{}, nil
	}

	var updatedTasks []interface{}
	for _, task := range data {
		taskMap, ok := task.(map[string]interface{})
		if !ok {
			continue
		}

		createdDate, createdOk := taskMap["createdDate"].(string)
		updatedDate, updatedOk := taskMap["updatedDate"].(string)

		if createdOk && updatedOk && createdDate != updatedDate {
			updatedTasks = append(updatedTasks, task)
		}
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
	return []interface{}{
		map[string]interface{}{
			"id":            "IEADTSKYA5CKABNW",
			"accountId":     "IEADTSKY",
			"title":         "Updated Task",
			"description":   "This task has been updated in Wrike",
			"status":        "Completed",
			"importance":    "High",
			"createdDate":   "2023-03-15T09:45:30.000Z",
			"updatedDate":   "2023-03-20T14:30:45.000Z",
			"completedDate": "2023-03-20T14:30:45.000Z",
		},
	}
}

func NewTaskUpdatedTrigger() sdk.Trigger {
	return &TaskUpdatedTrigger{}
}
