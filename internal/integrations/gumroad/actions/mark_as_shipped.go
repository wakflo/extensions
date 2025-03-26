package actions

import (
	"errors"

	"github.com/wakflo/extensions/internal/integrations/gumroad/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type markasShippedActionProps struct {
	SaleID string `json:"sale_id"`
}

type MarkasShippedAction struct{}

func (a *MarkasShippedAction) Name() string {
	return "Mark Product as Shipped"
}

func (a *MarkasShippedAction) Description() string {
	return "Mark a physical product as shipped."
}

func (a *MarkasShippedAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *MarkasShippedAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &markasShippedDocs,
	}
}

func (a *MarkasShippedAction) Icon() *string {
	icon := "mdi:package-variant-closed-multiple"
	return &icon
}

func (a *MarkasShippedAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"sale_id": shared.ListSalesInput("Sales ID", "Sales ID", true),
	}
}

func (a *MarkasShippedAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[markasShippedActionProps](ctx.BaseContext)
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

	sale, err := shared.DisableProduct(accessToken, input.SaleID)

	return sale, nil
}

func (a *MarkasShippedAction) Auth() *sdk.Auth {
	return nil
}

func (a *MarkasShippedAction) SampleData() sdkcore.JSON {
	return nil
}

func (a *MarkasShippedAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewMarkasShippedAction() sdk.Action {
	return &MarkasShippedAction{}
}
