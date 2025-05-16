// actions/get_product.go
package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/sendowl/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type getProductActionProps struct {
	ProductID string `json:"product_id"`
}

type GetProductAction struct{}

func (a *GetProductAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		DisplayName:   "Get Product",
		Description:   "Retrieve detailed information about a specific product from your SendOwl account using the product ID.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: getProductDocs,
		SampleOutput: map[string]any{
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
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *GetProductAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_product", "Get Product")

	shared.GetProductProp("product_id", "Product ID", "Select the ID of the product you want to retrieve.", true, form)

	schema := form.Build()

	return schema
}

func (a *GetProductAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getProductActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	url := "/products/" + input.ProductID

	product, err := shared.GetSendOwlClient(shared.BaseURL, authCtx.Extra["api_key"], authCtx.Extra["api_secret"], url)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (a *GetProductAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewGetProductAction() sdk.Action {
	return &GetProductAction{}
}
