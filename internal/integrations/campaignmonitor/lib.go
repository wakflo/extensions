package campaignmonitor

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/campaignmonitor/actions"
	"github.com/wakflo/extensions/internal/integrations/campaignmonitor/shared"
	"github.com/wakflo/extensions/internal/integrations/campaignmonitor/triggers"
	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewCampaignMonitor(), Flow, ReadME)

type CampaignMonitor struct{}

func (n *CampaignMonitor) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedAuth,
	}
}

func (n *CampaignMonitor) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewSubscriberAddedTrigger(),
		triggers.NewCampaignSentTrigger(),
	}
}

func (n *CampaignMonitor) Actions() []sdk.Action {
	return []sdk.Action{
		// actions.NewListCampaignsAction(),
		// actions.NewGetCampaignAction(),
		actions.NewCreateCampaignAction(),
		actions.NewListSubscribersAction(),
		actions.NewGetSubscriberListsAction(),
		actions.NewCreateCampaignFromTemplateAction(),
		actions.NewSendCampaignAction(),
		actions.NewListAllCampaignsAction(),
		// actions.NewAddSubscriberAction(),
	}
}

func NewCampaignMonitor() sdk.Integration {
	return &CampaignMonitor{}
}
