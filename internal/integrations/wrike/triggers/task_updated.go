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

type taskUpdatedTriggerProps struct {
	FolderID string `json:"folderId"`
	Limit    int    `json:"limit"`
}

type TaskUpdatedTrigger struct{}

func (t *TaskUpdatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "task_updated",
		DisplayName:   "Task Updated",
		Description:   "Triggers when a task is updated in your Wrike account.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: taskUpdatedDocs,
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

func (t *TaskUpdatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("task_updated", "Task Updated")

	shared.GetFoldersProp(form)

	form.NumberField("limit", "Limit").
		Placeholder("Enter limit").
		Required(false).
		HelpText("Maximum number of tasks to return when triggered (1-100).")

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
	input, err := sdk.InputToTypeSafely[taskUpdatedTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	lastRunTime, err := ctx.GetMetadata("lastrun")

	var updatedTime string
	if lastRunTime != nil {
		updatedTime = lastRunTime.(*time.Time).UTC().Format(time.RFC3339)
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

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken
	response, err := shared.GetWrikeClient(token, endpoint)
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

func (t *TaskUpdatedTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewTaskUpdatedTrigger() sdk.Trigger {
	return &TaskUpdatedTrigger{}
}
