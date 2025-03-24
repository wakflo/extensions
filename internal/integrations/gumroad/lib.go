package gumroad

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/gumroad/actions"
	"github.com/wakflo/extensions/internal/integrations/gumroad/shared"
	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewGumroad(), Flow, ReadME)

type Gumroad struct{}

func (n *Gumroad) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedAuth,
	}
}

func (n *Gumroad) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		// triggers.NewSaleCreatedTrigger(),
		// triggers.NewRefundCreatedTrigger(),
		// triggers.NewDisputeCreatedTrigger(),
		// triggers.NewSubscriptionCreatedTrigger(),
		// triggers.NewSubscriptionCancelledTrigger(),
	}
}

func (n *Gumroad) Actions() []sdk.Action {
	return []sdk.Action{
		// actions.NewGetProductAction(),
		actions.NewListProductsAction(),
		actions.NewGetProductAction(),
		actions.NewDisableProductAction(),
		actions.NewEnableProductAction(),
		actions.NewListSalesAction(),
		actions.NewGetSaleAction(),
		// actions.NewCreateProductAction(),
		// actions.NewUpdateProductAction(),
		// actions.NewGetSaleAction(),
		// actions.NewListSalesAction(),
		// actions.NewGetSubscriptionAction(),
		// actions.NewListSubscriptionsAction(),
		// actions.NewGetCustomerAction(),
		// actions.NewListCustomersAction(),
	}
}

func NewGumroad() sdk.Integration {
	return &Gumroad{}
}
