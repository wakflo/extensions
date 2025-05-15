package hubspot

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/hubspot/actions"
	"github.com/wakflo/extensions/internal/integrations/hubspot/shared"
	"github.com/wakflo/extensions/internal/integrations/hubspot/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewHubspot())

type Hubspot struct{}

func (n *Hubspot) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Hubspot) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.HubspotSharedAuth,
	}
}

func (n *Hubspot) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewContactUpdatedTrigger(),

		triggers.NewDealUpdatedTrigger(),

		triggers.NewTaskCreatedTrigger(),
	}
}

func (n *Hubspot) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewCreateContactAction(),

		actions.NewCreateTicketAction(),

		actions.NewListContactsAction(),

		actions.NewListTicketsAction(),

		actions.NewRetrieveContactAction(),

		actions.NewSearchOwnerByEmailAction(),

		actions.NewGetDealAction(),
	}
}

func NewHubspot() sdk.Integration {
	return &Hubspot{}
}
