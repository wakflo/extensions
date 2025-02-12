package actions

import (
	"context"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getProductActionProps struct {
	ProductID uint64 `json:"productId"`
}

type GetProductAction struct{}

func (a *GetProductAction) Name() string {
	return "Get Product"
}

func (a *GetProductAction) Description() string {
	return "Retrieves product information from the specified source, such as an e-commerce platform or inventory management system."
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
	return nil
}

func (a *GetProductAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"productId": autoform.NewNumberField().
			SetDisplayName("Product").
			SetDescription("The ID of the product.").
			SetRequired(true).
			Build(),
	}
}

func (a *GetProductAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getProductActionProps](ctx.BaseContext)
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
		return nil, fmt.Errorf("no product found with ID '%d'", input.ProductID)
	}
	productMap := map[string]interface{}{
		"ID":          product.Id,
		"Title":       product.Title,
		"Description": product.BodyHTML,
		"Price":       product.Variants[0].Price,
		"Variants":    product.Variants,
	}
	return sdk.JSON(map[string]interface{}{
		"product details": productMap,
	}), err
}

func (a *GetProductAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetProductAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *GetProductAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetProductAction() sdk.Action {
	return &GetProductAction{}
}
