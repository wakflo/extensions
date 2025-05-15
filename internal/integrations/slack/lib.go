package slack

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/slack/actions"
	"github.com/wakflo/extensions/internal/integrations/slack/shared"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewSlack())

type Slack struct{}

func (n *Slack) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Slack) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedSlackAuth,
	}
}

func (n *Slack) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *Slack) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewSendPublicChannelMessageAction(),

		actions.NewSendPrivateChannelMessageAction(),

		actions.NewSendDirectMessageAction(),
	}
}

func NewSlack() sdk.Integration {
	return &Slack{}
}
