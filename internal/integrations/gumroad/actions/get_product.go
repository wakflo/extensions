package actions

import (
	"errors"

	"github.com/wakflo/extensions/internal/integrations/gumroad/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getProductActionProps struct {
	ProductID string `json:"product_id"`
	// Category  *string `json:"category"`
}

type GetProductAction struct{}

func (a *GetProductAction) Name() string {
	return "Get Product"
}

func (a *GetProductAction) Description() string {
	return "Retrieves a product from your Gumroad store."
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
	icon := "mdi:package-variant-closed-multiple"
	return &icon
}

func (a *GetProductAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"product_id": shared.ListProductsInput("Product ID", "Product ID", true),
	}
}

func (a *GetProductAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getProductActionProps](ctx.BaseContext)
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

	product, err := shared.GetProduct(accessToken, input.ProductID)

	return product, nil
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
