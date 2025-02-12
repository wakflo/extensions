package triggers

import (
	"context"
	"fmt"
	"time"

	"github.com/hiscaler/woocommerce-go"
	"github.com/wakflo/extensions/internal/integrations/woocommerce/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type newOrderTriggerProps struct {
	Tag string `json:"tag"`
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
	return map[string]*sdkcore.AutoFormSchema{}
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
	_, err := sdk.InputToTypeSafely[newOrderTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	wooClient, err := shared.InitClient(ctx.BaseContext)
	if err != nil {
		return nil, err
	}
	lastRunTime := ctx.Metadata().LastRun

	var formattedTime string
	if lastRunTime != nil {
		utcTime := lastRunTime.UTC()
		formattedTime = utcTime.Format(time.RFC3339)
	}

	params := woocommerce.OrdersQueryParams{
		After: formattedTime,
	}

	newOrder, total, totalPages, isLastPage, err := wooClient.Services.Order.All(params)
	if err != nil {
		return nil, err
	}
	fmt.Println(totalPages, total, isLastPage)

	return newOrder, nil
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
