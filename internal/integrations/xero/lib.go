package xero

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/xero/actions"
	"github.com/wakflo/extensions/internal/integrations/xero/shared"
	"github.com/wakflo/extensions/internal/integrations/xero/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewXero())

type Xero struct{}

func (n *Xero) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Xero) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedAuth,
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
