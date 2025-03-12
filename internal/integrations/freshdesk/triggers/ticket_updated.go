package triggers

import (
	"context"
	"time"

	"github.com/wakflo/extensions/internal/integrations/freshdesk/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type TicketUpdatedTrigger struct{}

func (t *TicketUpdatedTrigger) Name() string {
	return "Ticket Updated"
}

func (t *TicketUpdatedTrigger) Description() string {
	return "Trigger a workflow when a ticket is updated in Freshdesk."
}

func (t *TicketUpdatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *TicketUpdatedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &ticketUpdatedDocs,
	}
}

func (t *TicketUpdatedTrigger) Icon() *string {
	icon := "mdi:ticket-edit"
	return &icon
}

func (t *TicketUpdatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{}
}

func (t *TicketUpdatedTrigger) Start(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *TicketUpdatedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *TicketUpdatedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	lastRunTime := ctx.Metadata().LastRun

	var updatedSince string
	if lastRunTime != nil {
		updatedSince = lastRunTime.UTC().Format(time.RFC3339)
	} else {
		updatedSince = ""
	}

	endpoint := "/tickets?order_by=updated_at&order_type=desc&updated_since=" + updatedSince

	domain := ctx.Auth.Extra["domain"]
	freshdeskDomain := "https://" + domain + ".freshdesk.com"

	response, err := shared.GetTickets(endpoint, freshdeskDomain, ctx.Auth.Extra["api-key"])
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (t *TicketUpdatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *TicketUpdatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *TicketUpdatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":           "123",
		"subject":      "Updated Support Request",
		"description":  "Updated help request details",
		"status":       "3",
		"priority":     "2",
		"requester_id": "456",
		"updated_at":   "2023-12-02T14:45:30Z",
	}
}

func NewTicketUpdatedTrigger() sdk.Trigger {
	return &TicketUpdatedTrigger{}
}
