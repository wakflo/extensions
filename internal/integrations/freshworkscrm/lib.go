package freshworkscrm

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/freshworkscrm/actions"
	"github.com/wakflo/extensions/internal/integrations/freshworkscrm/shared"
	"github.com/wakflo/extensions/internal/integrations/freshworkscrm/triggers"

	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewFreshWorksCRM())

type FreshWorksCRM struct{}

func (n *FreshWorksCRM) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *FreshWorksCRM) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.FreshworksCRMSharedAuth,
	}
}

func (n *FreshWorksCRM) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewContactCreatedTrigger(),
	}
}

func (n *FreshWorksCRM) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewCreateNewContactAction(),
		actions.NewUpdateContactAction(),
		actions.NewListContactsAction(),
	}
}

func NewFreshWorksCRM() sdk.Integration {
	return &FreshWorksCRM{}
}
