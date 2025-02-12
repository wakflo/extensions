package actions

import (
	"context"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getProductVariantActionProps struct {
	ProductID uint64 `json:"productId"`
	VariantID uint64 `json:"variantId"`
}

type GetProductVariantAction struct{}

func (a *GetProductVariantAction) Name() string {
	return "Get Product Variant"
}

func (a *GetProductVariantAction) Description() string {
	return "Retrieves product variant information based on provided input parameters, such as product ID or SKU. This action can be used to fetch details like variant name, price, and inventory levels, allowing you to incorporate this data into your workflow automation process."
}

func (a *GetProductVariantAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetProductVariantAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getProductVariantDocs,
	}
}

func (a *GetProductVariantAction) Icon() *string {
	return nil
}

func (a *GetProductVariantAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"productId": autoform.NewNumberField().
			SetDisplayName("Product ID").
			SetDescription("ID of product").
			SetRequired(false).
			Build(),
		"variantId": autoform.NewNumberField().
			SetDisplayName("Variant ID").
			SetDescription("product variant ID").
			SetRequired(false).
			Build(),
	}
}

func (a *GetProductVariantAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getProductVariantActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx.BaseContext)
	if err != nil {
		return nil, err
	}
	product, err := client.Product.Get(context.Background(), input.ProductID, nil)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, fmt.Errorf("no product variant found with ID '%d'", input.ProductID)
	}
	for _, variant := range product.Variants {
		if variant.Id == input.VariantID {
			return variant, nil
		}
	}
	return nil, nil
}

func (a *GetProductVariantAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetProductVariantAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *GetProductVariantAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetProductVariantAction() sdk.Action {
	return &GetProductVariantAction{}
}
