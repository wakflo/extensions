package actions

import (
	"github.com/wakflo/extensions/internal/integrations/woocommerce/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type findProductActionProps struct {
	ProductID int `json:"productId"`
}

type FindProductAction struct{}

func (a *FindProductAction) Name() string {
	return "Find Product"
}

func (a *FindProductAction) Description() string {
	return "Searches for a product by its name or ID in an external system, returning the matching product details."
}

func (a *FindProductAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *FindProductAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &findProductDocs,
	}
}

func (a *FindProductAction) Icon() *string {
	return nil
}

func (a *FindProductAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"productId": autoform.NewNumberField().
			SetDisplayName("Product ID").
			SetDescription("Enter product ID").
			SetRequired(true).
			Build(),
	}
}

func (a *FindProductAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[findProductActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	wooClient, err := shared.InitClient(ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	product, err := wooClient.Services.Product.One(input.ProductID)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (a *FindProductAction) Auth() *sdk.Auth {
	return nil
}

func (a *FindProductAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *FindProductAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewFindProductAction() sdk.Action {
	return &FindProductAction{}
}
