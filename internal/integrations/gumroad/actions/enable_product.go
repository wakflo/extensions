package actions

import (
	"errors"

	"github.com/wakflo/extensions/internal/integrations/gumroad/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type enableProductActionProps struct {
	ProductID string `json:"product_id"`
	// Category  *string `json:"category"`
}

type EnableProductAction struct{}

func (a *EnableProductAction) Name() string {
	return "Enable Product"
}

func (a *EnableProductAction) Description() string {
	return "Enable a product in your Gumroad store."
}

func (a *EnableProductAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *EnableProductAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getProductDocs,
	}
}

func (a *EnableProductAction) Icon() *string {
	icon := "mdi:package-variant-closed-multiple"
	return &icon
}

func (a *EnableProductAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"product_id": autoform.NewShortTextField().
			SetDisplayName("Product ID").
			SetDescription("Product ID").
			SetRequired(true).Build(),
	}
}

func (a *EnableProductAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[enableProductActionProps](ctx.BaseContext)
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

	product, err := shared.EnableProduct(accessToken, input.ProductID)

	return product, nil
}

func (a *EnableProductAction) Auth() *sdk.Auth {
	return nil
}

func (a *EnableProductAction) SampleData() sdkcore.JSON {
	return nil
}

func (a *EnableProductAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewEnableProductAction() sdk.Action {
	return &EnableProductAction{}
}
