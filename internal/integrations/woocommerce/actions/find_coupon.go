package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/woocommerce/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type findCouponActionProps struct {
	CouponID int `json:"couponId"`
}

type FindCouponAction struct{}

func (a *FindCouponAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "find_coupon",
		DisplayName:   "Find Coupon",
		Description:   "Searches for available coupons and discounts that can be applied to a specific order or transaction, allowing you to automate the process of finding the best deals and optimizing your customers' purchasing experiences.",
		Type:          core.ActionTypeAction,
		Documentation: findCouponDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *FindCouponAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("find_coupon", "Find Coupon")

	form.NumberField("couponId", "Coupon ID").
		Placeholder("Enter the coupon ID").
		Required(true).
		HelpText("Enter the coupon ID")

	schema := form.Build()

	return schema
}

func (a *FindCouponAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[findCouponActionProps](ctx)
	if err != nil {
		return nil, err
	}

	wooClient, err := shared.InitClient(ctx)
	if err != nil {
		return nil, err
	}

	coupon, err := wooClient.Services.Coupon.One(input.CouponID)
	if err != nil {
		return nil, err
	}

	return coupon, nil
}

func (a *FindCouponAction) Auth() *core.AuthMetadata {
	return nil
}

func NewFindCouponAction() sdk.Action {
	return &FindCouponAction{}
}
