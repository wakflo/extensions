package triggers

import (
	"context"
	"time"

	"github.com/wakflo/extensions/internal/integrations/freshdesk/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type TicketCreatedTrigger struct{}

func (t *TicketCreatedTrigger) Name() string {
	return "Ticket Created"
}

func (t *TicketCreatedTrigger) Description() string {
	return "Trigger a workflow when a new ticket is created in Freshdesk."
}

func (t *TicketCreatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *TicketCreatedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &ticketCreatedDocs,
	}
}

func (t *TicketCreatedTrigger) Icon() *string {
	icon := "mdi:ticket-plus"
	return &icon
}

func (t *TicketCreatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{}
}

// Start initializes the TicketCreatedTrigger
func (t *TicketCreatedTrigger) Start(ctx sdk.LifecycleContext) error {
	return nil
}

// Stop shuts down the TicketCreatedTrigger
func (t *TicketCreatedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

// Execute performs the polling for new tickets
func (t *TicketCreatedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	lastRunTime := ctx.Metadata().LastRun

	var createdSince string
	if lastRunTime != nil {
		createdSince = lastRunTime.UTC().Format(time.RFC3339)
	} else {
		createdSince = ""
	}

	endpoint := "/tickets?order_by=created_at&order_type=desc&created_since=" + createdSince
	domain := ctx.Auth.Extra["domain"]
	freshdeskDomain := "https://" + domain + ".freshdesk.com"

	response, err := shared.GetTickets(endpoint, freshdeskDomain, ctx.Auth.Extra["api-key"])
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (t *TicketCreatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *TicketCreatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *TicketCreatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":           "123",
		"subject":      "New Support Request",
		"description":  "I need help with your product",
		"status":       "2",
		"priority":     "1",
		"requester_id": "456",
		"created_at":   "2023-12-01T12:30:45Z",
	}
}

func NewTicketCreatedTrigger() sdk.Trigger {
	return &TicketCreatedTrigger{}
}
