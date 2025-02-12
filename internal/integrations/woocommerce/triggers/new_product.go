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

type newProductTriggerProps struct{}

type NewProductTrigger struct{}

func (t *NewProductTrigger) Name() string {
	return "New Product"
}

func (t *NewProductTrigger) Description() string {
	return "Triggered when a new product is created in your product information management system or e-commerce platform, allowing you to automate workflows and processes related to product launches, inventory management, and order fulfillment."
}

func (t *NewProductTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *NewProductTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &newProductDocs,
	}
}

func (t *NewProductTrigger) Icon() *string {
	return nil
}

func (t *NewProductTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{}
}

// Start initializes the newProductTrigger, required for event and webhook triggers in a lifecycle context.
func (t *NewProductTrigger) Start(ctx sdk.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the newProductTrigger, cleaning up resources and performing necessary teardown operations.
func (t *NewProductTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of newProductTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails. required for Pooling triggers
func (t *NewProductTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	_, err := sdk.InputToTypeSafely[newProductTriggerProps](ctx.BaseContext)
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

	params := woocommerce.ProductsQueryParams{
		After: formattedTime,
	}

	newProduct, total, totalPages, isLastPage, err := wooClient.Services.Product.All(params)
	if err != nil {
		return nil, err
	}
	fmt.Println(totalPages, total, isLastPage)

	return newProduct, nil
}

func (t *NewProductTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *NewProductTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *NewProductTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func NewNewProductTrigger() sdk.Trigger {
	return &NewProductTrigger{}
}
