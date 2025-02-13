package stripe

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/stripe/actions"
	"github.com/wakflo/extensions/internal/integrations/stripe/shared"
	"github.com/wakflo/extensions/internal/integrations/stripe/triggers"
	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewStripe(), Flow, ReadME)

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
