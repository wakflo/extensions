package triggers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/hubspot/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type taskCreatedTriggerProps struct {
	Properties string `json:"properties"`
}

type TaskCreatedTrigger struct{}

func (t *TaskCreatedTrigger) Name() string {
	return "Task Created"
}

func (t *TaskCreatedTrigger) Description() string {
	return "Trigger a workflow when new tasks are created in your HubSpot CRM"
}

func (t *TaskCreatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *TaskCreatedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &taskCreatedDoc,
	}
}

func (t *TaskCreatedTrigger) Icon() *string {
	icon := "mdi:calendar-check"
	return &icon
}

func (t *TaskCreatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"properties": autoform.NewShortTextField().
			SetDisplayName("Task Properties").
			SetDescription("Comma-separated list of properties to retrieve (e.g., hs_task_subject,hs_task_body,hs_task_priority)").
			SetRequired(false).
			Build(),
	}
}

func (t *TaskCreatedTrigger) Start(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *TaskCreatedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *TaskCreatedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	props, err := sdk.InputToTypeSafely[taskCreatedTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	lastRunTime := ctx.Metadata().LastRun
	url := "/crm/v3/objects/tasks/search"

	const limit = 100

	requestBody := map[string]interface{}{
		"limit": limit,
		"sorts": []map[string]string{
			{
				"propertyName": "hs_createdate",
				"direction":    "DESCENDING",
			},
		},
	}

	if lastRunTime != nil {
		requestBody["filterGroups"] = []map[string]interface{}{
			{
				"filters": []map[string]interface{}{
					{
						"propertyName": "hs_createdate",
						"operator":     "GT",
						"value":        lastRunTime.UnixMilli(),
					},
				},
			},
		}
	}

	if props.Properties != "" {
		requestBody["properties"] = append(
			[]string{"hs_task_subject", "hs_task_body", "hs_task_priority", "hs_createdate"},
			props.Properties,
		)
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	resp, err := shared.HubspotClient(url, ctx.Auth.AccessToken, http.MethodPost, jsonBody)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *TaskCreatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *TaskCreatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *TaskCreatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"results": []map[string]any{},
	}
}

func NewTaskCreatedTrigger() sdk.Trigger {
	return &TaskCreatedTrigger{}
}
