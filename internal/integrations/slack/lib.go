package slack

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/slack/actions"
	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewSlack(), Flow, ReadME)

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
