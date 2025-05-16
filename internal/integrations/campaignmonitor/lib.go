package campaignmonitor

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/campaignmonitor/actions"
	"github.com/wakflo/extensions/internal/integrations/campaignmonitor/shared"
	"github.com/wakflo/extensions/internal/integrations/campaignmonitor/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewCampaignMonitor())

type CampaignMonitor struct{}

func (n *CampaignMonitor) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *CampaignMonitor) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.CampaignMonitorSharedAuth,
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
		actions.NewCreateCampaignAction(),
		actions.NewListSubscribersAction(),
		actions.NewGetSubscriberListsAction(),
		actions.NewCreateCampaignFromTemplateAction(),
		actions.NewSendCampaignAction(),
		actions.NewListAllCampaignsAction(),
		actions.NewAddSubscriberAction(),
		actions.NewGetCampaignListsAndSegmentsAction(),
	}
}

func NewCampaignMonitor() sdk.Integration {
	return &CampaignMonitor{}
}
