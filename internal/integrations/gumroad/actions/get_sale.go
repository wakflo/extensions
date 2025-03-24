package actions

import (
	"errors"

	"github.com/wakflo/extensions/internal/integrations/gumroad/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getSaleActionProps struct {
	SaleID string `json:"sale_id"`
}

type GetSaleAction struct{}

func (a *GetSaleAction) Name() string {
	return "Get Sale"
}

func (a *GetSaleAction) Description() string {
	return "Retrieves a sale from your Gumroad store."
}

func (a *GetSaleAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetSaleAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getSaleDocs,
	}
}

func (a *GetSaleAction) Icon() *string {
	icon := "mdi:package-variant-closed-multiple"
	return &icon
}

func (a *GetSaleAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"sale_id": autoform.NewShortTextField().
			SetDisplayName("Sale ID").
			SetDescription("Sale ID").
			SetRequired(false).Build(),
	}
}

func (a *GetSaleAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getSaleActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing Gumroad auth token")
	}
	accessToken := ctx.Auth.AccessToken

	if input.SaleID == "" {
		return nil, errors.New("sale ID is required")
	}

	sale, err := shared.GetSale(accessToken, input.SaleID)

	return sale, nil
}

func (a *GetSaleAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetSaleAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *GetSaleAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetSaleAction() sdk.Action {
	return &GetSaleAction{}
}
