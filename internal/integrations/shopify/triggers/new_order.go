package triggers

import (
	"context"
	"time"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type newOrderTriggerProps struct {
	Name string `json:"name"`
}

type NewOrderTrigger struct{}

func (t *NewOrderTrigger) Name() string {
	return "New Order"
}

func (t *NewOrderTrigger) Description() string {
	return "Triggered when a new order is created in your e-commerce platform or inventory management system, allowing you to automate tasks and workflows immediately after an order is placed."
}

func (t *NewOrderTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *NewOrderTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &newOrderDocs,
	}
}

func (t *NewOrderTrigger) Icon() *string {
	return nil
}

func (t *NewOrderTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"name": autoform.NewShortTextField().
			SetLabel("Name").
			SetRequired(true).
			SetPlaceholder("Your name").
			Build(),
	}
}

// Start initializes the newOrderTrigger, required for event and webhook triggers in a lifecycle context.
func (t *NewOrderTrigger) Start(ctx sdk.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the newOrderTrigger, cleaning up resources and performing necessary teardown operations.
func (t *NewOrderTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of newOrderTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails. required for Pooling triggers
func (t *NewOrderTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
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

	orders, err := client.Order.List(context.Background(), options)
	if err != nil {
		return nil, err
	}

	return sdk.JSON(map[string]interface{}{
		"order details": orders,
	}), err
}

func (t *NewOrderTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *NewOrderTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *NewOrderTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func NewNewOrderTrigger() sdk.Trigger {
	return &NewOrderTrigger{}
}
