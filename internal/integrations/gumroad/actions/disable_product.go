package actions

import (
	"errors"

	"github.com/wakflo/extensions/internal/integrations/gumroad/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type disableProductActionProps struct {
	ProductID string `json:"product_id"`
}

type DisableProductAction struct{}

func (a *DisableProductAction) Name() string {
	return "Disable Product"
}

func (a *DisableProductAction) Description() string {
	return "Disable a product in your Gumroad store."
}

func (a *DisableProductAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *DisableProductAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getProductDocs,
	}
}

func (a *DisableProductAction) Icon() *string {
	icon := "mdi:package-variant-closed-multiple"
	return &icon
}

func (a *DisableProductAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"product_id": shared.ListProductsInput("Product ID", "Product ID", true),
	}
}

func (a *DisableProductAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[disableProductActionProps](ctx.BaseContext)
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

	product, err := shared.DisableProduct(accessToken, input.ProductID)

	return product, nil
}

func (a *DisableProductAction) Auth() *sdk.Auth {
	return nil
}

func (a *DisableProductAction) SampleData() sdkcore.JSON {
	return nil
}

func (a *DisableProductAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewDisableProductAction() sdk.Action {
	return &DisableProductAction{}
}
