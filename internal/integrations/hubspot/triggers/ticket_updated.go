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

type ticketCreatedTriggerProps struct {
	Properties string `json:"properties"`
}

type TicketCreatedTrigger struct{}

func (t *TicketCreatedTrigger) Name() string {
	return "Ticket Created"
}

func (t *TicketCreatedTrigger) Description() string {
	return "Trigger a workflow when new tickets are created in your HubSpot CRM"
}

func (t *TicketCreatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *TicketCreatedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &ticketCreatedDoc,
	}
}

func (t *TicketCreatedTrigger) Icon() *string {
	icon := "mdi:ticket-confirmation"
	return &icon
}

func (t *TicketCreatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"properties": autoform.NewShortTextField().
			SetDisplayName("Ticket Properties").
			SetDescription("Comma-separated list of properties to retrieve (e.g., subject,hs_ticket_priority,hs_pipeline_stage)").
			SetRequired(false).
			Build(),
	}
}

func (t *TicketCreatedTrigger) Start(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *TicketCreatedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *TicketCreatedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	props, err := sdk.InputToTypeSafely[ticketCreatedTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	lastRunTime := ctx.Metadata().LastRun
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
	resp, err := shared.HubspotClient(url, ctx.Auth.AccessToken, http.MethodPost, jsonBody)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *TicketCreatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *TicketCreatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *TicketCreatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"results": []map[string]any{},
	}
}

func NewTicketCreatedTrigger() sdk.Trigger {
	return &TicketCreatedTrigger{}
}
