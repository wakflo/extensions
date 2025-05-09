// triggers/product_updated.go
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

type productUpdatedTriggerProps struct{}

type ProductUpdatedTrigger struct{}

func (t *ProductUpdatedTrigger) Name() string {
	return "Product Updated"
}

func (t *ProductUpdatedTrigger) Description() string {
	return "Automatically trigger workflows when a product in your SendOwl account is updated. This allows you to maintain synchronized product information across systems or notify team members about product changes."
}

func (t *ProductUpdatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *ProductUpdatedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &productUpdatedDocs,
	}
}

func (t *ProductUpdatedTrigger) Icon() *string {
	icon := "mdi:package-variant-closed-check"
	return &icon
}

func (t *ProductUpdatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{}
}

// Start initializes the ProductUpdatedTrigger, required for event and webhook triggers in a lifecycle context.
func (t *ProductUpdatedTrigger) Start(ctx sdk.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the ProductUpdatedTrigger, cleaning up resources and performing necessary teardown operations.
func (t *ProductUpdatedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of ProductUpdatedTrigger by processing the input context and returning a JSON response.
// It retrieves products updated since the last check.
func (t *ProductUpdatedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	lastRunTime := ctx.Metadata().LastRun

	url := "/products"

	response, err := shared.GetSendOwlClient(shared.AltBaseURL, ctx.Auth.Extra["api_key"], ctx.Auth.Extra["api_secret"], url)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error fetching products: %v", err))
	}

	// Check if the response is an array as expected
	if !response.IsArray {
		return nil, errors.New("unexpected response format: expected array of orders")
	}

	// Extract products from the response
	var products []map[string]interface{}
	for _, item := range response.Array {
		// Each item should be a map with an "product" field
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue // skip if not a map
		}

		// Get the product object
		productRaw, ok := itemMap["product"]
		if !ok {
			continue // skip if no product field
		}

		// Convert the product to a map
		productMap, ok := productRaw.(map[string]interface{})
		if !ok {
			continue // skip if product isn't a map
		}

		products = append(products, productMap)
	}

	// If no last run time, return all products
	if lastRunTime == nil {
		return products, nil
	}

	// Filter orders based on updated_at time
	var filteredProducts []map[string]interface{}
	for _, product := range products {
		updatedAtStr, ok := product["updated_at"].(string)
		if !ok {
			continue // skip products without updated_at
		}

		updatedAt, err := time.Parse(time.RFC3339, updatedAtStr)
		if err != nil {
			continue // skip if timestamp is invalid
		}

		if updatedAt.After(*lastRunTime) {
			filteredProducts = append(filteredProducts, product)
		}
	}

	return filteredProducts, nil
}

func (t *ProductUpdatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *ProductUpdatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *ProductUpdatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":              123456,
		"name":            "Digital Marketing Guide",
		"price":           29.99,
		"currency":        "USD",
		"product_type":    "pdf",
		"description":     "A comprehensive guide to digital marketing strategies and tactics.",
		"stock_level":     999,
		"sales_count":     156,
		"download_limit":  5,
		"download_expiry": 30,
		"created_at":      "2023-07-15T10:30:00Z",
		"updated_at":      "2023-08-20T14:45:22Z",
	}
}

func NewProductUpdatedTrigger() sdk.Trigger {
	return &ProductUpdatedTrigger{}
}
