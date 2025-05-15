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

type getProductActionProps struct {
	ProductID uint64 `json:"productId"`
}

type GetProductAction struct{}

func (a *GetProductAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_product",
		DisplayName:   "Get Product",
		Description:   "Retrieves product information from the specified source, such as an e-commerce platform or inventory management system.",
		Type:          core.ActionTypeAction,
		Documentation: getProductDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *GetProductAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_product", "Get Product")

	form.NumberField("productId", "Product").
		Placeholder("The ID of the product.").
		Required(true).
		HelpText("The ID of the product.")

	schema := form.Build()
	return schema
}

func (a *GetProductAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getProductActionProps](ctx)
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
		return nil, fmt.Errorf("no product found with ID '%d'", input.ProductID)
	}
	productMap := map[string]interface{}{
		"ID":          product.Id,
		"Title":       product.Title,
		"Description": product.BodyHTML,
		"Price":       product.Variants[0].Price,
		"Variants":    product.Variants,
	}
	return core.JSON(map[string]interface{}{
		"product details": productMap,
	}), err
}

func (a *GetProductAction) Auth() *core.AuthMetadata {
	return nil
}

func NewGetProductAction() sdk.Action {
	return &GetProductAction{}
}
