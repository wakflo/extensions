package discord

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/discord/actions"
	"github.com/wakflo/extensions/internal/integrations/discord/shared"
	"github.com/wakflo/extensions/internal/integrations/discord/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewDiscord())

type Discord struct{}

func (d *Discord) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (d *Discord) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedAuth,
	}
}

func (d *Discord) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewMessageReceivedTrigger(),
	}
}

func (d *Discord) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewSendMessageAction(),
		actions.NewCreateChannelAction(),
		actions.NewFindChannelAction(),
		actions.NewRenameChannelAction(),
		actions.NewDeleteChannelAction(),
		actions.NewBanGuildMemberAction(),
		actions.NewRemoveGuildMemberAction(),
		actions.NewUnbanGuildMemberAction(),
		actions.NewListGuildMembersAction(),
		actions.NewAddRoleToMemberAction(),
	}
}

func NewDiscord() sdk.Integration {
	return &Discord{}
}
