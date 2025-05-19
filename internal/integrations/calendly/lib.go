package calendly

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/calendly/actions"
	"github.com/wakflo/extensions/internal/integrations/calendly/shared"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewCalendly())

type Calendly struct{}

func (n *Calendly) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Calendly) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedCalendlyAuth,
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
