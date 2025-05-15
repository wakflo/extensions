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

type newOrderTriggerProps struct {
	CreatedTime *time.Time `json:"createdTime"`
}

type NewOrderTrigger struct{}

func (t *NewOrderTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "new_order",
		DisplayName:   "New Order",
		Description:   "Triggered when a new order is created in your Shopify store, allowing you to automate tasks and workflows immediately after an order is placed.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: newOrderDocs,
		SampleOutput: map[string]any{
			"orders": []map[string]any{
				{
					"id":                 123456789,
					"name":               "#1001",
					"email":              "customer@example.com",
					"created_at":         "2023-01-01T12:00:00Z",
					"financial_status":   "paid",
					"fulfillment_status": "fulfilled",
					"total_price":        "125.00",
				},
			},
		},
	}
}

func (t *NewOrderTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func (t *NewOrderTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *NewOrderTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("shopify-new-order", "New Order")
	schema := form.Build()

	return schema
}

// Start initializes the newOrderTrigger, required for event and webhook triggers in a lifecycle context.
func (t *NewOrderTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the newOrderTrigger, cleaning up resources and performing necessary teardown operations.
func (t *NewOrderTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of newOrderTrigger by processing the input context and returning a JSON response.
func (t *NewOrderTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[newOrderTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx)
	if err != nil {
		return nil, err
	}

	var lastRunTime time.Time
	if input.CreatedTime == nil {
		lr, err := ctx.GetMetadata("lastRun")
		if err == nil && lr != nil {
			lastRunTime = *lr.(*time.Time)
		}
	} else {
		lastRunTime = *input.CreatedTime
	}

	options := &goshopify.ListOptions{
		CreatedAtMin: lastRunTime.UTC(),
	}

	// Get orders from Shopify
	orders, err := client.Order.List(context.Background(), options)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (t *NewOrderTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *NewOrderTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":                 123456789,
		"name":               "#1001",
		"email":              "customer@example.com",
		"created_at":         "2023-01-01T12:00:00Z",
		"financial_status":   "paid",
		"fulfillment_status": "fulfilled",
		"total_price":        "125.00",
	}
}

func NewNewOrderTrigger() sdk.Trigger {
	return &NewOrderTrigger{}
}
