package freshworkscrm

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/freshworkscrm/actions"
	"github.com/wakflo/extensions/internal/integrations/freshworkscrm/shared"
	"github.com/wakflo/extensions/internal/integrations/freshworkscrm/triggers"

	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewFreshWorksCRM(), Flow, ReadME)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

type FreshWorksCRM struct{}

func (n *FreshWorksCRM) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.FreshworksSharedAuth,
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
