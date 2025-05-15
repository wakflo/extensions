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

type contactUpdatedTriggerProps struct {
	Properties string `json:"properties"`
}

type ContactUpdatedTrigger struct{}

// Metadata returns metadata about the trigger
func (t *ContactUpdatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "contact_updated",
		DisplayName:   "Contact Updated",
		Description:   "Trigger a workflow when a contact is created or updated in your HubSpot CRM",
		Type:          core.TriggerTypePolling,
		Documentation: contactUpdatedDoc,
		Icon:          "mdi:account-check",
		SampleOutput: map[string]any{
			"results": []map[string]any{
				{
					"id": "51",
					"properties": map[string]any{
						"firstname":        "John",
						"lastname":         "Doe",
						"email":            "john.doe@example.com",
						"phone":            "+1234567890",
						"createdate":       "2023-03-15T09:31:40.678Z",
						"lastmodifieddate": "2023-03-15T10:45:12.412Z",
					},
					"createdAt": "2023-03-15T09:31:40.678Z",
					"updatedAt": "2023-03-15T10:45:12.412Z",
				},
			},
		},
	}
}

// Auth returns the authentication requirements for the trigger
func (t *ContactUpdatedTrigger) Auth() *core.AuthMetadata {
	return nil
}

// GetType returns the type of the trigger
func (t *ContactUpdatedTrigger) GetType() core.TriggerType {
	return core.TriggerTypePolling
}

// Props returns the schema for the trigger's input configuration
func (t *ContactUpdatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("contact_updated", "Contact Updated")

	form.TextareaField("properties", "Contact Properties").
		Required(false).
		HelpText("Comma-separated list of contact properties to include in the response (e.g., firstname,lastname,email)")

	schema := form.Build()

	return schema
}

// Start initializes the trigger, required for event and webhook triggers in a lifecycle context
func (t *ContactUpdatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Stop shuts down the trigger, cleaning up resources and performing necessary teardown operations
func (t *ContactUpdatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the trigger logic
func (t *ContactUpdatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[contactUpdatedTriggerProps](ctx)
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

	url := "/crm/v3/objects/contacts/search"
	const limit = 100

	requestBody := map[string]interface{}{
		"limit": limit,
		"sorts": []map[string]string{
			{
				"propertyName": "lastmodifieddate",
				"direction":    "DESCENDING",
			},
		},
	}

	if lastRunTime != nil {
		requestBody["filterGroups"] = []map[string]interface{}{
			{
				"filters": []map[string]interface{}{
					{
						"propertyName": "lastmodifieddate",
						"operator":     "GT",
						"value":        lastRunTime.UnixMilli(),
					},
				},
			},
		}
	}

	if input.Properties != "" {
		requestBody["properties"] = append(
			[]string{"firstname", "lastname", "email", "lastmodifieddate"},
			input.Properties,
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
func (t *ContactUpdatedTrigger) Criteria(ctx context.Context) core.TriggerCriteria {
	return core.TriggerCriteria{}
}

// SampleData returns sample data for this trigger
func (t *ContactUpdatedTrigger) SampleData() core.JSON {
	return map[string]any{
		"results": []map[string]any{
			{
				"id": "51",
				"properties": map[string]any{
					"firstname":        "John",
					"lastname":         "Doe",
					"email":            "john.doe@example.com",
					"phone":            "+1234567890",
					"createdate":       "2023-03-15T09:31:40.678Z",
					"lastmodifieddate": "2023-03-15T10:45:12.412Z",
				},
				"createdAt": "2023-03-15T09:31:40.678Z",
				"updatedAt": "2023-03-15T10:45:12.412Z",
			},
		},
	}
}

func NewContactUpdatedTrigger() sdk.Trigger {
	return &ContactUpdatedTrigger{}
}
