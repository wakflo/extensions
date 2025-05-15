package harvest

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/harvest/actions"
	"github.com/wakflo/extensions/internal/integrations/harvest/shared"
	"github.com/wakflo/extensions/internal/integrations/harvest/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewHarvest())

type Harvest struct{}

func (n *Harvest) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Harvest) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedHarvestAuth,
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
