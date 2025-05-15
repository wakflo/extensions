package shopify

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/shopify/actions"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/extensions/internal/integrations/shopify/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewShopify())

type Shopify struct{}

func (n *Shopify) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Shopify) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.ShopifySharedAuth,
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
