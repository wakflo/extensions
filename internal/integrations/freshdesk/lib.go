package freshdesk

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/freshdesk/actions"
	"github.com/wakflo/extensions/internal/integrations/freshdesk/shared"
	"github.com/wakflo/extensions/internal/integrations/freshdesk/triggers"

	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewFreshdesk())

type Freshdesk struct{}

func (n *Freshdesk) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Freshdesk) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.FreshdeskSharedAuth,
	}
}

func (n *Freshdesk) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewTicketCreatedTrigger(),
		triggers.NewTicketUpdatedTrigger(),
	}
}

func (n *Freshdesk) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewCreateTicketAction(),
		actions.NewGetTicketAction(),
		actions.NewListTicketsAction(),
		actions.NewUpdateTicketAction(),
		actions.NewSearchTicketsAction(),
	}
}

func NewFreshdesk() sdk.Integration {
	return &Freshdesk{}
}
