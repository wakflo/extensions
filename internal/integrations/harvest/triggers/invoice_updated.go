package triggers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/harvest/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type invoiceUpdatedTriggerProps struct{}

type InvoiceUpdatedTrigger struct{}

// Metadata returns metadata about the trigger
func (t *InvoiceUpdatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "invoice_updated",
		DisplayName:   "Invoice Updated",
		Description:   "Invoice Updated integration trigger is designed to automate workflows when an invoice is updated in your accounting or ERP system. This trigger can be used to initiate a series of automated tasks, such as sending notifications to stakeholders, updating CRM records, or triggering payment processing. When an invoice is updated, this trigger will fire and allow you to define the subsequent actions that need to take place.",
		Type:          core.TriggerTypePolling,
		Documentation: invoiceUpdatedDocs,
		SampleOutput: []map[string]any{
			{
				"id":         13150453,
				"client_id":  5735776,
				"number":     "1001",
				"amount":     288.23,
				"due_amount": 288.23,
				"subject":    "Web Design",
				"state":      "open",
				"issue_date": "2023-05-01",
				"due_date":   "2023-05-31",
				"sent_at":    "2023-05-02T14:30:00Z",
				"updated_at": "2023-05-10T09:45:22Z",
			},
		},
	}
}

// Auth returns the authentication requirements for the trigger
func (t *InvoiceUpdatedTrigger) Auth() *core.AuthMetadata {
	return nil
}

// GetType returns the type of the trigger
func (t *InvoiceUpdatedTrigger) GetType() core.TriggerType {
	return core.TriggerTypePolling
}

// Props returns the schema for the trigger's input configuration
func (t *InvoiceUpdatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("invoice_updated", "Invoice Updated")

	// The original doesn't have any properties

	schema := form.Build()

	return schema
}

// Start initializes the trigger, required for event and webhook triggers in a lifecycle context
func (t *InvoiceUpdatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Stop shuts down the trigger, cleaning up resources and performing necessary teardown operations
func (t *InvoiceUpdatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the trigger logic
func (t *InvoiceUpdatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (core.JSON, error) {
	// Get the last run time
	var lastRunTime *time.Time
	lr, err := ctx.GetMetadata("lastRun")
	if err == nil && lr != nil {
		lastRunTime = lr.(*time.Time)
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	var updatedTime string
	if lastRunTime != nil {
		updatedTime = lastRunTime.UTC().Format(time.RFC3339)
	}

	url := "/v2/invoices?updated_since=" + updatedTime

	response, err := shared.GetHarvestClient(authCtx.Token.AccessToken, url)
	if err != nil {
		return nil, fmt.Errorf("Error fetching data: %v", err)
	}

	data, ok := response["invoices"].([]interface{})
	if !ok {
		return nil, errors.New("invalid response format: data field is not an array")
	}

	return data, nil
}

// Criteria returns the criteria for triggering this trigger
func (t *InvoiceUpdatedTrigger) Criteria(ctx context.Context) core.TriggerCriteria {
	return core.TriggerCriteria{}
}

// SampleData returns sample data for this trigger
func (t *InvoiceUpdatedTrigger) SampleData() core.JSON {
	return []map[string]any{
		{
			"id":         13150453,
			"client_id":  5735776,
			"number":     "1001",
			"amount":     288.23,
			"due_amount": 288.23,
			"subject":    "Web Design",
			"state":      "open",
			"issue_date": "2023-05-01",
			"due_date":   "2023-05-31",
			"sent_at":    "2023-05-02T14:30:00Z",
			"updated_at": "2023-05-10T09:45:22Z",
		},
	}
}

func NewInvoiceUpdatedTrigger() sdk.Trigger {
	return &InvoiceUpdatedTrigger{}
}
