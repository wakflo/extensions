package actions

import (
	"github.com/wakflo/extensions/internal/integrations/woocommerce/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type findCouponActionProps struct {
	CouponID int `json:"couponId"`
}

type FindCouponAction struct{}

func (a *FindCouponAction) Name() string {
	return "Find Coupon"
}

func (a *FindCouponAction) Description() string {
	return "Searches for available coupons and discounts that can be applied to a specific order or transaction, allowing you to automate the process of finding the best deals and optimizing your customers' purchasing experiences."
}

func (a *FindCouponAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *FindCouponAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &findCouponDocs,
	}
}

func (a *FindCouponAction) Icon() *string {
	return nil
}

func (a *FindCouponAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"couponId": autoform.NewNumberField().
			SetDisplayName("Coupon ID").
			SetDescription("Enter the coupon ID").
			SetRequired(true).
			Build(),
	}
}

func (a *FindCouponAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[findCouponActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	wooClient, err := shared.InitClient(ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	product, err := wooClient.Services.Coupon.One(input.CouponID)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (a *FindCouponAction) Auth() *sdk.Auth {
	return nil
}

func (a *FindCouponAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *FindCouponAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewFindCouponAction() sdk.Action {
	return &FindCouponAction{}
}
