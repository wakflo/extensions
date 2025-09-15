package telegrambot

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/telegrambot/actions"
	"github.com/wakflo/extensions/internal/integrations/telegrambot/shared"
	"github.com/wakflo/extensions/internal/integrations/telegrambot/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewTelegram())

type Telegram struct{}

func (t *Telegram) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (t *Telegram) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedAuth,
	}
}

func (t *Telegram) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewMessageReceivedTrigger(),
	}
}

func (t *Telegram) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewSendMessageAction(),
		actions.NewSendPhotoAction(),
		actions.NewCreateInviteLinkAction(),
		actions.NewGetChatAdministratorsAction(),
		actions.NewGetChatMemberCountAction(),
		actions.NewGetChatMemberAction(),
	}
}

func NewTelegram() sdk.Integration {
	return &Telegram{}
}
