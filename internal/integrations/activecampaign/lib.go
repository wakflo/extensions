package activecampaign

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/activecampaign/actions"
	"github.com/wakflo/extensions/internal/integrations/activecampaign/shared"
	"github.com/wakflo/extensions/internal/integrations/activecampaign/triggers"
	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewActiveCampaign(), Flow, ReadME)

type ActiveCampaign struct{}

func (a *ActiveCampaign) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedAuth,
	}
}

func (a *ActiveCampaign) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewContactUpdatedTrigger(),
		triggers.NewContactCreatedTrigger(),
	}
}

func (a *ActiveCampaign) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewListContactsAction(),
		actions.NewGetContactAction(),
		actions.NewCreateContactAction(),
		actions.NewUpdateContactAction(),
	}
}

func NewActiveCampaign() sdk.Integration {
	return &ActiveCampaign{}
}
