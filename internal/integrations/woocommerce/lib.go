package woocommerce

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/woocommerce/actions"
	"github.com/wakflo/extensions/internal/integrations/woocommerce/shared"
	"github.com/wakflo/extensions/internal/integrations/woocommerce/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewWooCommerce())

type WooCommerce struct{}

func (n *WooCommerce) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *WooCommerce) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.WooSharedAuth,
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

		actions.NewGetCustomerByIDAction(),

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
