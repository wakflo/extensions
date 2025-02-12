package slack

import (
	"github.com/wakflo/extensions/internal/integrations/slack/actions"
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewSlack())

type Slack struct{}

func (n *Slack) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: false,
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
