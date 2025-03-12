package freshdesk

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/freshdesk/actions"
	"github.com/wakflo/extensions/internal/integrations/freshdesk/shared"
	// "github.com/wakflo/extensions/internal/integrations/freshdesk/triggers"

	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewFreshdesk(), Flow, ReadME)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

type Freshdesk struct{}

func (n *Freshdesk) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedAuth,
	}
}

func (n *Freshdesk) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *Freshdesk) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewCreateTicketAction(),
	}
}

func NewFreshdesk() sdk.Integration {
	return &Freshdesk{}
}
