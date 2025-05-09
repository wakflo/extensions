package triggers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/wakflo/extensions/internal/integrations/sendowl/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type orderCompletedTriggerProps struct{}

type OrderCompletedTrigger struct{}

func (t *OrderCompletedTrigger) Name() string {
	return "Order Completed"
}

func (t *OrderCompletedTrigger) Description() string {
	return "Automatically trigger workflows when a new order is marked as completed in your SendOwl account. This allows you to automate post-purchase processes such as customer onboarding, fulfillment, or marketing follow-ups."
}

func (t *OrderCompletedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *OrderCompletedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &orderCompletedDocs,
	}
}

func (t *OrderCompletedTrigger) Icon() *string {
	icon := "mdi:cart-check"
	return &icon
}

func (t *OrderCompletedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{}
}

// Start initializes the OrderCompletedTrigger, required for event and webhook triggers in a lifecycle context.
func (t *OrderCompletedTrigger) Start(ctx sdk.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the OrderCompletedTrigger, cleaning up resources and performing necessary teardown operations.
func (t *OrderCompletedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *OrderCompletedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	lastRunTime := ctx.Metadata().LastRun

	// Build query to get orders completed since last run
	url := "/orders?status=completed"

	if lastRunTime != nil {
		formattedTime := lastRunTime.UTC().Format("2006-01-02")
		url = url + "&from_date=" + formattedTime
	}

	response, err := shared.GetSendOwlClient(shared.AltBaseURL, ctx.Auth.Extra["api_key"], ctx.Auth.Extra["api_secret"], url)
	if err != nil {
		return nil, fmt.Errorf("error fetching completed orders: %v", err)
	}

	// Check if the response is an array as expected
	if !response.IsArray {
		return nil, errors.New("unexpected response format: expected array of orders")
	}

	// Extract orders from the response
	var orders []map[string]interface{}
	for _, item := range response.Array {
		// Each item should be a map with an "order" field
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue // skip if not a map
		}

		// Get the order object
		orderRaw, ok := itemMap["order"]
		if !ok {
			continue // skip if no order field
		}

		// Convert the order to a map
		orderMap, ok := orderRaw.(map[string]interface{})
		if !ok {
			continue // skip if order isn't a map
		}

		orders = append(orders, orderMap)
	}

	// If no last run time, return all orders (assuming they're all completed)
	if lastRunTime == nil {
		return orders, nil
	}

	// Filter orders based on completed_at time
	var filteredOrders []map[string]interface{}
	for _, order := range orders {
		completedAtStr, ok := order["completed_at"].(string)
		if !ok {
			continue // skip orders without completed_at
		}

		completedAt, err := time.Parse(time.RFC3339, completedAtStr)
		if err != nil {
			continue // skip if timestamp is invalid
		}

		if completedAt.After(*lastRunTime) {
			filteredOrders = append(filteredOrders, order)
		}
	}

	return filteredOrders, nil
}

func (t *OrderCompletedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *OrderCompletedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *OrderCompletedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":             123456,
		"buyer_name":     "John Doe",
		"buyer_email":    "john.doe@example.com",
		"total":          29.99,
		"currency":       "USD",
		"status":         "completed",
		"payment_method": "credit_card",
		"transaction_id": "tx_abc123def456",
		"created_at":     "2023-08-15T14:30:22Z",
		"completed_at":   "2023-08-15T14:35:45Z",
		"products": []map[string]any{
			{
				"id":           789012,
				"name":         "Digital Marketing Guide",
				"price":        29.99,
				"product_type": "pdf",
			},
		},
	}
}

func NewOrderCompletedTrigger() sdk.Trigger {
	return &OrderCompletedTrigger{}
}
