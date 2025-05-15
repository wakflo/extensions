package triggers

import (
	"context"
	"time"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type newCustomerTriggerProps struct{}

type NewCustomerTrigger struct{}

func (t *NewCustomerTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "new_customer",
		DisplayName:   "New Customer",
		Description:   "Triggered when a new customer is created in your Shopify store, this integration allows you to automate workflows and tasks immediately after a new customer is added, streamlining your sales and marketing processes.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: newCustomerDocs,
		SampleOutput: map[string]any{
			"customers": []map[string]any{
				{
					"id":           123456789,
					"email":        "customer@example.com",
					"first_name":   "John",
					"last_name":    "Doe",
					"phone":        "+1234567890",
					"created_at":   "2023-01-01T12:00:00Z",
					"updated_at":   "2023-01-01T12:00:00Z",
					"orders_count": 5,
				},
			},
		},
	}
}

func (t *NewCustomerTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func (t *NewCustomerTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *NewCustomerTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("shopify-new-customer", "New Customer")

	form.TextField("email", "Email Filter").
		Placeholder("john@example.com").
		Required(false).
		HelpText("Only trigger for customers with this email address. Leave blank to trigger for all new customers.")

	schema := form.Build()

	return schema
}

// Start initializes the newCustomerTrigger, required for event and webhook triggers in a lifecycle context.
func (t *NewCustomerTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the newCustomerTrigger, cleaning up resources and performing necessary teardown operations.
func (t *NewCustomerTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of newCustomerTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
func (t *NewCustomerTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	client, err := shared.CreateClient(ctx)
	if err != nil {
		return nil, err
	}

	// Get the last run time
	var lastRunTime time.Time
	lr, err := ctx.GetMetadata("lastRun")
	if err == nil && lr != nil {
		lastRunTime = *lr.(*time.Time)
	}

	options := &goshopify.ListOptions{
		CreatedAtMin: lastRunTime.UTC(),
	}

	// Get customers from Shopify
	customers, err := client.Customer.List(context.Background(), options)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (t *NewCustomerTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *NewCustomerTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":           123456789,
		"email":        "customer@example.com",
		"first_name":   "John",
		"last_name":    "Doe",
		"phone":        "+1234567890",
		"created_at":   "2023-01-01T12:00:00Z",
		"updated_at":   "2023-01-01T12:00:00Z",
		"orders_count": 5,
	}
}

func NewNewCustomerTrigger() sdk.Trigger {
	return &NewCustomerTrigger{}
}
