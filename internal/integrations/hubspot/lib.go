package hubspot

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/hubspot/actions"
	"github.com/wakflo/extensions/internal/integrations/hubspot/shared"
	"github.com/wakflo/extensions/internal/integrations/hubspot/triggers"

	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewHubspot(), Flow, ReadME)

type Hubspot struct{}

func (n *Hubspot) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.HubspotSharedAuth,
	}
}

func (n *Hubspot) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewContactUpdatedTrigger(),

		triggers.NewDealUpdatedTrigger(),
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
