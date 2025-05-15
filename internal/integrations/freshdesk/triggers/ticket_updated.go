package triggers

import (
	"context"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/freshdesk/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type TicketUpdatedTrigger struct{}

func (t *TicketUpdatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "ticket_updated",
		DisplayName:   "Ticket Updated",
		Description:   "Trigger a workflow when a ticket is updated in Freshdesk.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: ticketUpdatedDocs,
		Icon:          "mdi:ticket-edit",
		SampleOutput: map[string]any{
			"id":           "123",
			"subject":      "Updated Support Request",
			"description":  "Updated help request details",
			"status":       "3",
			"priority":     "2",
			"requester_id": "456",
			"updated_at":   "2023-12-02T14:45:30Z",
		},
	}
}

func (t *TicketUpdatedTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func (t *TicketUpdatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *TicketUpdatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("freshdesk-ticket-updated", "Ticket Updated")

	// No properties needed for this trigger

	schema := form.Build()

	return schema
}

// Start initializes the TicketUpdatedTrigger
func (t *TicketUpdatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Stop shuts down the TicketUpdatedTrigger
func (t *TicketUpdatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the polling for updated tickets
func (t *TicketUpdatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	var lastRunTime *time.Time
	lastRun, err := ctx.GetMetadata("lastRun")
	if err == nil && lastRun != nil {
		lastRunTime = lastRun.(*time.Time)
	}

	var updatedSince string
	if lastRunTime != nil {
		updatedSince = lastRunTime.UTC().Format(time.RFC3339)
	} else {
		updatedSince = ""
	}

	endpoint := "/tickets?order_by=updated_at&order_type=desc&updated_since=" + updatedSince

	domain := authCtx.Extra["domain"]
	freshdeskDomain := "https://" + domain + ".freshdesk.com"

	response, err := shared.GetTickets(endpoint, freshdeskDomain, authCtx.Extra["api-key"])
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (t *TicketUpdatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func NewTicketUpdatedTrigger() sdk.Trigger {
	return &TicketUpdatedTrigger{}
}
