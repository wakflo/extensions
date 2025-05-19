package triggers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/zendeskapp/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type ticketCreatedProps struct {
	// ID string `json:"id"`
}

type TicketCreatedTrigger struct{}

func (t *TicketCreatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "new_ticket_created",
		DisplayName:   "New Ticket Added",
		Description:   "Triggers workflow when a new ticket is added",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: newInvoiceDocs,
		SampleOutput: map[string]interface{}{
			"message": "New ticket created",
		},
	}
}

func (t *TicketCreatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("new_ticket_created", "New Ticket Added")

	schema := form.Build()

	return schema
}

// Start initializes the TicketCreatedTrigger, required for event and webhook triggers in a lifecycle context.
func (t *TicketCreatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the TicketCreatedTrigger, cleaning up resources and performing necessary teardown operations.
func (t *TicketCreatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of TicketCreatedTrigger by processing the input context and returning a JSON response.
// It checks for new tickets in Zendesk and returns them.
func (t *TicketCreatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	authData := ctx.Auth()
	if authData == nil {
		return nil, errors.New("missing authentication data")
	}

	// Check if required auth details are present
	email, ok := authData.Extra["email"]
	if !ok || email == "" {
		return nil, errors.New("missing zendesk email")
	}

	apiToken, ok := authData.Extra["api-token"]
	if !ok || apiToken == "" {
		return nil, errors.New("missing zendesk api-token")
	}

	subdomain, ok := authData.Extra["subdomain"]
	if !ok || subdomain == "" {
		return nil, errors.New("missing zendesk subdomain")
	}

	baseURL := fmt.Sprintf("https://%s.zendesk.com/api/v2", subdomain)
	fullURL := baseURL + "/search.json?query=type:ticket"

	lastRunTime, err := ctx.GetMetadata("lastrun")
	if lastRunTime != nil {
		createdAfter := lastRunTime.(*time.Time).UTC().Format("2006-01-02T15:04:05Z")
		fullURL = fmt.Sprintf("%s+created>=%s", fullURL, createdAfter)
	}

	// Using the existing zendeskRequest function from the original code
	response, err := shared.ZendeskRequest(http.MethodGet, fullURL, email, apiToken, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching tickets: %v", err)
	}

	return response, nil
}

func (t *TicketCreatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *TicketCreatedTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func (t *TicketCreatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "New ticket created in Zendesk!",
	}
}

func NewTicketCreatedTrigger() sdk.Trigger {
	return &TicketCreatedTrigger{}
}
