package harvest

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/harvest/actions"
	"github.com/wakflo/extensions/internal/integrations/harvest/shared"
	"github.com/wakflo/extensions/internal/integrations/harvest/triggers"
	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewHarvest(), Flow, ReadME)

type Harvest struct{}

func (n *Harvest) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedAuth,
	}
}

func (n *Harvest) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewInvoiceUpdatedTrigger(),
	}
}

func (n *Harvest) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewListInvoicesAction(),

		actions.NewGetInvoiceAction(),
	}
}

func NewHarvest() sdk.Integration {
	return &Harvest{}
}
