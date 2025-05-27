package triggers

import (
	"context"
	"fmt"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/freshdesk/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type TicketCreatedTrigger struct{}

func (t *TicketCreatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "ticket_created",
		DisplayName:   "Ticket Created",
		Description:   "Trigger a workflow when a new ticket is created in Freshdesk.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: ticketCreatedDocs,
		Icon:          "mdi:ticket-plus",
		SampleOutput: map[string]any{
			"id":           "123",
			"subject":      "New Support Request",
			"description":  "I need help with your product",
			"status":       "2",
			"priority":     "1",
			"requester_id": "456",
			"created_at":   "2023-12-01T12:30:45Z",
		},
	}
}

func (t *TicketCreatedTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func (t *TicketCreatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *TicketCreatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("freshdesk-ticket-created", "Ticket Created")

	// No properties needed for this trigger
	schema := form.Build()

	return schema
}

// Start initializes the TicketCreatedTrigger
func (t *TicketCreatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Stop shuts down the TicketCreatedTrigger
func (t *TicketCreatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the polling for new tickets
func (t *TicketCreatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	// Get last run time from metadata
	var lastRunTime *time.Time
	lastRun, err := ctx.GetMetadata("lastRun")
	if err == nil && lastRun != nil {
		lastRunTime = lastRun.(*time.Time)
	}

	var endpoint string
	if lastRunTime != nil {
		updatedSince := lastRunTime.UTC().Format(time.RFC3339)
		endpoint = fmt.Sprintf("/tickets?order_by=created_at&order_type=desc&updated_since=%s", updatedSince)
	} else {
		// If no lastRunTime, just get recent tickets
		endpoint = "/tickets?order_by=created_at&order_type=desc"
	}

	domain := authCtx.Extra["domain"]
	freshdeskDomain := "https://" + domain + ".freshdesk.com"

	response, err := shared.GetTickets(endpoint, freshdeskDomain, authCtx.Extra["api-key"])
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (t *TicketCreatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func NewTicketCreatedTrigger() sdk.Trigger {
	return &TicketCreatedTrigger{}
}
