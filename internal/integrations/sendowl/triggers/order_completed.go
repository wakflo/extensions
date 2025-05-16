package triggers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/sendowl/shared"

	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type orderCompletedTriggerProps struct{}

type OrderCompletedTrigger struct{}

func (t *OrderCompletedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "order_completed",
		DisplayName:   "Order Completed",
		Description:   "Triggers a workflow when a new order is marked as completed in your SendOwl account.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: orderCompletedDocs,
		SampleOutput: map[string]any{
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
		},
	}
}

func (t *OrderCompletedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *OrderCompletedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("order_completed", "Order Completed")

	schema := form.Build()

	return schema
}

// Start initializes the OrderCompletedTrigger, required for event and webhook triggers in a lifecycle context.
func (t *OrderCompletedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the OrderCompletedTrigger, cleaning up resources and performing necessary teardown operations.
func (t *OrderCompletedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *OrderCompletedTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	lr, err := ctx.GetMetadata("lastRun")
	if err != nil {
		return nil, err
	}

	lastRunTime := lr.(*time.Time)

	// Build query to get orders completed since last run
	url := "/orders?status=completed"

	if lastRunTime != nil {
		formattedTime := lastRunTime.UTC().Format("2006-01-02")
		url = url + "&from_date=" + formattedTime
	}

	response, err := shared.GetSendOwlClient(shared.AltBaseURL, authCtx.Extra["api_key"], authCtx.Extra["api_secret"], url)
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

func (t *OrderCompletedTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewOrderCompletedTrigger() sdk.Trigger {
	return &OrderCompletedTrigger{}
}
