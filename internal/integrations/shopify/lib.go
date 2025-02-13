package shopify

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/shopify/actions"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/extensions/internal/integrations/shopify/triggers"
	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewShopify(), Flow, ReadME)

type Shopify struct{}

func (n *Shopify) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedAuth,
	}
}

func (n *Shopify) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewNewOrderTrigger(),

		triggers.NewNewCustomerTrigger(),
	}
}

func (n *Shopify) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewUpdateCustomerAction(),

		actions.NewCloseOrderAction(),

		actions.NewUpdateProductAction(),

		actions.NewUpdateOrderAction(),

		actions.NewListProductsAction(),

		actions.NewListOrdersAction(),

		actions.NewListDraftOrdersAction(),

		actions.NewListCustomersAction(),

		actions.NewGetTransactionAction(),

		actions.NewGetProductVariantAction(),

		actions.NewGetProductAction(),

		actions.NewGetOrderAction(),

		actions.NewGetLocationsAction(),

		actions.NewGetCustomerOrderAction(),

		actions.NewGetCustomerAction(),

		actions.NewCreateTransactionAction(),

		actions.NewCreateProductAction(),

		actions.NewCreateOrderAction(),

		actions.NewCreateDraftOrderAction(),

		actions.NewCreateCustomerAction(),

		actions.NewCreateCollectAction(),

		actions.NewCancelOrderAction(),

		actions.NewAdjustInventoryLevelAction(),
	}
}

func NewShopify() sdk.Integration {
	return &Shopify{}
}
