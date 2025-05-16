package activecampaign

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/activecampaign/actions"
	"github.com/wakflo/extensions/internal/integrations/activecampaign/shared"
	"github.com/wakflo/extensions/internal/integrations/activecampaign/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewActiveCampaign())

type ActiveCampaign struct{}

func (a *ActiveCampaign) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (a *ActiveCampaign) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.ActiveCampaignSharedAuth,
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
