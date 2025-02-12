package woocommerce

import (
	"github.com/wakflo/extensions/internal/integrations/woocommerce/actions"
	"github.com/wakflo/extensions/internal/integrations/woocommerce/shared"
	"github.com/wakflo/extensions/internal/integrations/woocommerce/triggers"
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewWooCommerce())

type WooCommerce struct{}

func (n *WooCommerce) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedAuth,
	}
}

func (n *WooCommerce) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewNewProductTrigger(),

		triggers.NewNewOrderTrigger(),
	}
}

func (n *WooCommerce) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewUpdateProductAction(),

		actions.NewListProductsAction(),

		actions.NewListOrdersAction(),

		actions.NewGetCustomerByIdAction(),

		actions.NewFindProductAction(),

		actions.NewFindCustomerAction(),

		actions.NewFindCouponAction(),

		actions.NewCreateProductAction(),

		actions.NewCreateCustomerAction(),
	}
}

func NewWooCommerce() sdk.Integration {
	return &WooCommerce{}
}
