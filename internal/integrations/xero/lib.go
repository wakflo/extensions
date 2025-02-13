package xero

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/xero/actions"
	"github.com/wakflo/extensions/internal/integrations/xero/shared"
	"github.com/wakflo/extensions/internal/integrations/xero/triggers"
	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewXero(), Flow, ReadME)

type Xero struct{}

func (n *Xero) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedAuth,
	}
}

func (n *Xero) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewNewInvoiceTrigger(),
	}
}

func (n *Xero) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewGetInvoiceAction(),

		actions.NewListInvoicesAction(),

		actions.NewEmailInvoiceAction(),

		actions.NewCreateInvoiceAction(),
	}
}

func NewXero() sdk.Integration {
	return &Xero{}
}
