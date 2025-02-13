package calendly

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/calendly/actions"
	"github.com/wakflo/extensions/internal/integrations/calendly/shared"
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewCalendly(), Flow, ReadME)

type Calendly struct{}

func (n *Calendly) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: false,
		Schema:   *shared.SharedAuth,
	}
}

func (n *Calendly) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *Calendly) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewCreateSingleUseScheduleLinkAction(),
		actions.NewListEventsAction(),
		actions.NewGetEventAction(),
	}
}

func NewCalendly() sdk.Integration {
	return &Calendly{}
}

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string
