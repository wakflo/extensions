package stripe

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/stripe/actions"
	"github.com/wakflo/extensions/internal/integrations/stripe/shared"
	"github.com/wakflo/extensions/internal/integrations/stripe/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewStripe())

type Stripe struct{}

func (n *Stripe) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Stripe) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.StripeSharedAuth,
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

		actions.NewRetrieveCustomerAction(),
	}
}

func NewStripe() sdk.Integration {
	return &Stripe{}
}
