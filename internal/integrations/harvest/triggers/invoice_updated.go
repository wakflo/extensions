package triggers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/wakflo/extensions/internal/integrations/harvest/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type invoiceUpdatedTriggerProps struct{}

type InvoiceUpdatedTrigger struct{}

func (t *InvoiceUpdatedTrigger) Name() string {
	return "Invoice Updated"
}

func (t *InvoiceUpdatedTrigger) Description() string {
	return "Invoice Updated integration trigger is designed to automate workflows when an invoice is updated in your accounting or ERP system. This trigger can be used to initiate a series of automated tasks, such as sending notifications to stakeholders, updating CRM records, or triggering payment processing. When an invoice is updated, this trigger will fire and allow you to define the subsequent actions that need to take place."
}

func (t *InvoiceUpdatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *InvoiceUpdatedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &invoiceUpdatedDocs,
	}
}

func (t *InvoiceUpdatedTrigger) Icon() *string {
	return nil
}

func (t *InvoiceUpdatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{}
}

// Start initializes the invoiceUpdatedTrigger, required for event and webhook triggers in a lifecycle context.
func (t *InvoiceUpdatedTrigger) Start(ctx sdk.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the invoiceUpdatedTrigger, cleaning up resources and performing necessary teardown operations.
func (t *InvoiceUpdatedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of invoiceUpdatedTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails. required for Pooling triggers
func (t *InvoiceUpdatedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	lastRunTime := ctx.Metadata().LastRun

	var updatedTime string
	if lastRunTime != nil {
		updatedTime = lastRunTime.UTC().Format(time.RFC3339)
	}

	url := "/v2/invoices?updated_since=" + updatedTime

	response, err := shared.GetHarvestClient(ctx.Auth.AccessToken, url)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error fetching data: %v", err))
	}

	data, ok := response["invoices"].([]interface{})
	if !ok {
		return nil, errors.New("invalid response format: data field is not an array")
	}

	return data, nil
}

func (t *InvoiceUpdatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *InvoiceUpdatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *InvoiceUpdatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func NewInvoiceUpdatedTrigger() sdk.Trigger {
	return &InvoiceUpdatedTrigger{}
}
