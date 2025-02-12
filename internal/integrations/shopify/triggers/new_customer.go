package triggers

import (
	"context"
	"time"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type newCustomerTriggerProps struct {
	Name string `json:"name"`
}

type NewCustomerTrigger struct{}

func (t *NewCustomerTrigger) Name() string {
	return "New Customer"
}

func (t *NewCustomerTrigger) Description() string {
	return "Triggered when a new customer is created in your CRM or database, this integration allows you to automate workflows and tasks immediately after a new customer is added, streamlining your sales and marketing processes."
}

func (t *NewCustomerTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *NewCustomerTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &newCustomerDocs,
	}
}

func (t *NewCustomerTrigger) Icon() *string {
	return nil
}

func (t *NewCustomerTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{}
}

// Start initializes the newCustomerTrigger, required for event and webhook triggers in a lifecycle context.
func (t *NewCustomerTrigger) Start(ctx sdk.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the newCustomerTrigger, cleaning up resources and performing necessary teardown operations.
func (t *NewCustomerTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of newCustomerTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails. required for Pooling triggers
func (t *NewCustomerTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	client, err := shared.CreateClient(ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	var lastRunTime time.Time
	if ctx.Metadata().LastRun != nil {
		lastRunTime = *ctx.Metadata().LastRun
	}

	options := &goshopify.ListOptions{
		CreatedAtMin: lastRunTime.UTC(),
	}

	customers, err := client.Customer.List(context.Background(), options)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (t *NewCustomerTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *NewCustomerTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *NewCustomerTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func NewNewCustomerTrigger() sdk.Trigger {
	return &NewCustomerTrigger{}
}
