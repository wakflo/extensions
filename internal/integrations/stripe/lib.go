package stripe

import (
	"github.com/wakflo/extensions/internal/integrations/stripe/actions"
	"github.com/wakflo/extensions/internal/integrations/stripe/shared"
	"github.com/wakflo/extensions/internal/integrations/stripe/triggers"
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewStripe())

type Stripe struct{}

func (n *Stripe) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedAuth,
	}
}

func (n *Stripe) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewNewCustomerTrigger(),
	}
}

func (n *Stripe) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewSearchCustomerAction(),

		actions.NewCreateInvoiceAction(),

		actions.NewCreateCustomerAction(),
	}
}

func NewStripe() sdk.Integration {
	return &Stripe{}
}
