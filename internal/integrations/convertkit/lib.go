package convertkit

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/convertkit/actions"
	"github.com/wakflo/extensions/internal/integrations/convertkit/shared"
	"github.com/wakflo/extensions/internal/integrations/convertkit/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewConvertKit())

type ConvertKit struct{}

func (n *ConvertKit) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *ConvertKit) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.ConvertKitSharedAuth,
	}
}

func (n *ConvertKit) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewSubscriberCreatedTrigger(),
		triggers.NewTagCreatedTrigger(),
	}
}

func (n *ConvertKit) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewListSubscribersAction(),
		actions.NewGetSubscriberAction(),
		// actions.NewCreateSubscriberAction(),
		actions.NewCreateTagAction(),
		actions.NewTagSubscriberAction(),
		actions.NewListTagsAction(),
	}
}

func NewConvertKit() sdk.Integration {
	return &ConvertKit{}
}
