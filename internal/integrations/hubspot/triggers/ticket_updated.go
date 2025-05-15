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

type ticketCreatedTriggerProps struct {
	Properties string `json:"properties"`
}

type TicketCreatedTrigger struct{}

// Metadata returns metadata about the trigger
func (t *TicketCreatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "ticket_created",
		DisplayName:   "Ticket Created",
		Description:   "Trigger a workflow when new tickets are created in your HubSpot CRM",
		Type:          core.TriggerTypePolling,
		Documentation: ticketCreatedDoc,
		Icon:          "mdi:ticket-confirmation",
		SampleOutput: map[string]any{
			"results": []map[string]any{
				{
					"id": "12345",
					"properties": map[string]any{
						"subject":            "Technical issue with product",
						"content":            "User is experiencing login problems with the application.",
						"hs_ticket_priority": "HIGH",
						"hs_pipeline":        "0",
						"hs_pipeline_stage":  "1",
						"createdate":         "2023-04-15T09:30:00Z",
					},
					"createdAt": "2023-04-15T09:30:00Z",
					"updatedAt": "2023-04-15T09:30:00Z",
				},
			},
		},
	}
}

// Auth returns the authentication requirements for the trigger
func (t *TicketCreatedTrigger) Auth() *core.AuthMetadata {
	return nil
}

// GetType returns the type of the trigger
func (t *TicketCreatedTrigger) GetType() core.TriggerType {
	return core.TriggerTypePolling
}

// Props returns the schema for the trigger's input configuration
func (t *TicketCreatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("ticket_created", "Ticket Created")

	form.TextareaField("properties", "Ticket Properties").
		Required(false).
		HelpText("Comma-separated list of properties to retrieve (e.g., subject,hs_ticket_priority,hs_pipeline_stage)")

	schema := form.Build()

	return schema
}

// Start initializes the trigger, required for event and webhook triggers in a lifecycle context
func (t *TicketCreatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Stop shuts down the trigger, cleaning up resources and performing necessary teardown operations
func (t *TicketCreatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the trigger logic
func (t *TicketCreatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (core.JSON, error) {
	props, err := sdk.InputToTypeSafely[ticketCreatedTriggerProps](ctx)
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

	url := "/crm/v3/objects/tickets/search"
	const limit = 100

	requestBody := map[string]interface{}{
		"limit": limit,
		"sorts": []map[string]string{
			{
				"propertyName": "createdate",
				"direction":    "DESCENDING",
			},
		},
	}

	// Only add filter if there's a last run time
	if lastRunTime != nil {
		requestBody["filterGroups"] = []map[string]interface{}{
			{
				"filters": []map[string]interface{}{
					{
						"propertyName": "createdate",
						"operator":     "GT",
						"value":        lastRunTime.UnixMilli(),
					},
				},
			},
		}
	}

	// Add properties if specified
	if props.Properties != "" {
		requestBody["properties"] = append(
			[]string{"subject", "hs_ticket_priority", "hs_pipeline_stage", "createdate"},
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
func (t *TicketCreatedTrigger) Criteria(ctx context.Context) core.TriggerCriteria {
	return core.TriggerCriteria{}
}

// SampleData returns sample data for this trigger
func (t *TicketCreatedTrigger) SampleData() core.JSON {
	return map[string]any{
		"results": []map[string]any{
			{
				"id": "12345",
				"properties": map[string]any{
					"subject":            "Technical issue with product",
					"content":            "User is experiencing login problems with the application.",
					"hs_ticket_priority": "HIGH",
					"hs_pipeline":        "0",
					"hs_pipeline_stage":  "1",
					"createdate":         "2023-04-15T09:30:00Z",
				},
				"createdAt": "2023-04-15T09:30:00Z",
				"updatedAt": "2023-04-15T09:30:00Z",
			},
		},
	}
}

func NewTicketCreatedTrigger() sdk.Trigger {
	return &TicketCreatedTrigger{}
}
