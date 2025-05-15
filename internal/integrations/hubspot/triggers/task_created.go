package triggers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/hubspot/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type taskCreatedTriggerProps struct {
	Properties string `json:"properties"`
}

type TaskCreatedTrigger struct{}

// Metadata returns metadata about the trigger
func (t *TaskCreatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "task_created",
		DisplayName:   "Task Created",
		Description:   "Trigger a workflow when new tasks are created in your HubSpot CRM",
		Type:          core.TriggerTypePolling,
		Documentation: taskCreatedDoc,
		Icon:          "mdi:calendar-check",
		SampleOutput: map[string]any{
			"results": []map[string]any{
				{
					"id": "12345",
					"properties": map[string]any{
						"hs_task_subject":  "Follow up with customer",
						"hs_task_body":     "Discuss upcoming renewal and potential upsell opportunities",
						"hs_task_priority": "HIGH",
						"hs_task_status":   "NOT_STARTED",
						"hs_task_type":     "CALL",
						"hs_createdate":    "2023-04-15T09:30:00Z",
						"hs_timestamp":     "2023-04-20T14:00:00Z",
					},
					"createdAt": "2023-04-15T09:30:00Z",
					"updatedAt": "2023-04-15T09:30:00Z",
				},
			},
		},
	}
}

// Auth returns the authentication requirements for the trigger
func (t *TaskCreatedTrigger) Auth() *core.AuthMetadata {
	return nil
}

// GetType returns the type of the trigger
func (t *TaskCreatedTrigger) GetType() core.TriggerType {
	return core.TriggerTypePolling
}

// Props returns the schema for the trigger's input configuration
func (t *TaskCreatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("task_created", "Task Created")

	form.TextareaField("properties", "Task Properties").
		Required(false).
		HelpText("Comma-separated list of properties to retrieve (e.g., hs_task_subject,hs_task_body,hs_task_priority)")

	schema := form.Build()

	return schema
}

// Start initializes the trigger, required for event and webhook triggers in a lifecycle context
func (t *TaskCreatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Stop shuts down the trigger, cleaning up resources and performing necessary teardown operations
func (t *TaskCreatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the trigger logic
func (t *TaskCreatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (core.JSON, error) {
	props, err := sdk.InputToTypeSafely[taskCreatedTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	// Get the last run time
	var lastRunTime *time.Time
	lr, err := ctx.GetMetadata("lastRun")
	if err == nil && lr != nil {
		lastRunTime = lr.(*time.Time)
	}

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

	resp, err := shared.HubspotClient(url, authCtx.Token.AccessToken, http.MethodPost, jsonBody)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Criteria returns the criteria for triggering this trigger
func (t *TaskCreatedTrigger) Criteria(ctx context.Context) core.TriggerCriteria {
	return core.TriggerCriteria{}
}

// SampleData returns sample data for this trigger
func (t *TaskCreatedTrigger) SampleData() core.JSON {
	return map[string]any{
		"results": []map[string]any{
			{
				"id": "12345",
				"properties": map[string]any{
					"hs_task_subject":  "Follow up with customer",
					"hs_task_body":     "Discuss upcoming renewal and potential upsell opportunities",
					"hs_task_priority": "HIGH",
					"hs_task_status":   "NOT_STARTED",
					"hs_task_type":     "CALL",
					"hs_createdate":    "2023-04-15T09:30:00Z",
					"hs_timestamp":     "2023-04-20T14:00:00Z",
				},
				"createdAt": "2023-04-15T09:30:00Z",
				"updatedAt": "2023-04-15T09:30:00Z",
			},
		},
	}
}

func NewTaskCreatedTrigger() sdk.Trigger {
	return &TaskCreatedTrigger{}
}
