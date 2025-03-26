package actions

import (
	"errors"

	"github.com/wakflo/extensions/internal/integrations/gumroad/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type deleteProductActionProps struct {
	ProductID string `json:"product_id"`
}

type DeleteProductAction struct{}

func (a *DeleteProductAction) Name() string {
	return "Delete Product"
}

func (a *DeleteProductAction) Description() string {
	return "Permanently a product from your Gumroad store."
}

func (a *DeleteProductAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *DeleteProductAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &deleteProductDocs,
	}
}

func (a *DeleteProductAction) Icon() *string {
	icon := "mdi:package-variant-closed-multiple"
	return &icon
}

func (a *DeleteProductAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"product_id": shared.ListProductsInput("Product ID", "Product ID", true),
	}
}

func (a *DeleteProductAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[deleteProductActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing Gumroad auth token")
	}
	accessToken := ctx.Auth.AccessToken

	if input.ProductID == "" {
		return nil, errors.New("product ID is required")
	}

	product, err := shared.DeleteProduct(accessToken, input.ProductID)

	return product, nil
}

func (a *DeleteProductAction) Auth() *sdk.Auth {
	return nil
}

func (a *DeleteProductAction) SampleData() sdkcore.JSON {
	return nil
}

func (a *DeleteProductAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewDeleteProductAction() sdk.Action {
	return &DeleteProductAction{}
}
