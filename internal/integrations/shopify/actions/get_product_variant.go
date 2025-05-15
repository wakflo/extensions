package actions

import (
	"context"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getProductVariantActionProps struct {
	ProductID uint64 `json:"productId"`
	VariantID uint64 `json:"variantId"`
}

type GetProductVariantAction struct{}

func (a *GetProductVariantAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_product_variant",
		DisplayName:   "Get Product Variant",
		Description:   "Retrieves product variant information based on provided input parameters, such as product ID or SKU. This action can be used to fetch details like variant name, price, and inventory levels, allowing you to incorporate this data into your workflow automation process.",
		Type:          core.ActionTypeAction,
		Documentation: getProductVariantDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *GetProductVariantAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_product_variant", "Get Product Variant")

	form.NumberField("productId", "Product ID").
		Placeholder("ID of product").
		HelpText("ID of product")

	form.NumberField("variantId", "Variant ID").
		Placeholder("product variant ID").
		HelpText("product variant ID")

	schema := form.Build()
	return schema
}

func (a *GetProductVariantAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getProductVariantActionProps](ctx)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx)
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

func (a *GetProductVariantAction) Auth() *core.AuthMetadata {
	return nil
}

func NewGetProductVariantAction() sdk.Action {
	return &GetProductVariantAction{}
}
