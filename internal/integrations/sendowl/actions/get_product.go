// actions/get_product.go
package actions

import (
	"github.com/wakflo/extensions/internal/integrations/sendowl/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getProductActionProps struct {
	ProductID string `json:"product_id"`
}

type GetProductAction struct{}

func (a *GetProductAction) Name() string {
	return "Get Product"
}

func (a *GetProductAction) Description() string {
	return "Retrieve detailed information about a specific product from your SendOwl account using the product ID."
}

func (a *GetProductAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetProductAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getProductDocs,
	}
}

func (a *GetProductAction) Icon() *string {
	icon := "mdi:package-variant"
	return &icon
}

func (a *GetProductAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"product_id": shared.GetProductInput("Product ID", "Enter the ID of the product you want to retrieve.", true),
	}
}

func (a *GetProductAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getProductActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	url := "/products/" + input.ProductID

	product, err := shared.GetSendOwlClient(ctx.Auth.Extra["api_key"], ctx.Auth.Extra["api_secret"], url)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (a *GetProductAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetProductAction) SampleData() sdkcore.JSON {
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

func (a *GetProductAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetProductAction() sdk.Action {
	return &GetProductAction{}
}
