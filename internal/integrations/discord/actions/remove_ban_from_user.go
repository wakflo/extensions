package actions

import (
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/discord/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type unbanGuildMemberActionProps struct {
	GuildID string `json:"guild-id"`
	UserID  string `json:"user-id"`
	Reason  string `json:"reason"`
}

type UnbanGuildMemberAction struct{}

// Metadata returns metadata about the action
func (a *UnbanGuildMemberAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "unban_guild_member",
		DisplayName:   "Unban Guild Member",
		Description:   "Remove a ban from a user, allowing them to rejoin the Discord server",
		Type:          core.ActionTypeAction,
		Documentation: unbanGuildMemberDocs,
		SampleOutput: map[string]any{
			"success":  true,
			"user_id":  "123456789012345678",
			"guild_id": "857347647235678912",
			"action":   "unbanned",
			"reason":   "Appeal approved",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *UnbanGuildMemberAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("unban_guild_member", "Unban Guild Member")

	shared.RegisterGuildsInput(form, "Guilds", "List of guilds", true)

	form.TextField("user-id", "User ID").
		Required(true).
		HelpText("The ID of the user to unban from the server")

	form.TextField("reason", "Reason").
		Required(false).
		HelpText("Optional reason for unbanning the user (will appear in audit log)")

	schema := form.Build()
	return schema
}

func (a *UnbanGuildMemberAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[unbanGuildMemberActionProps](ctx)
	if err != nil {
		return nil, err
	}

	endpoint := "/guilds/" + input.GuildID + "/bans/" + input.UserID

	if input.Reason != "" {
		endpoint += "?reason=" + input.Reason
	}

	response, err := shared.GetDiscordClient(ctx.Auth().Key, endpoint, "DELETE", nil)
	if err != nil {
		return nil, err
	}

	fmt.Println(response...)

	return map[string]interface{}{
		"success":  true,
		"user_id":  input.UserID,
		"guild_id": input.GuildID,
		"action":   "unbanned",
		"reason":   input.Reason,
		"message":  "User successfully unbanned from guild",
	}, nil
}

func (a *UnbanGuildMemberAction) Auth() *core.AuthMetadata {
	return nil
}

func NewUnbanGuildMemberAction() sdk.Action {
	return &UnbanGuildMemberAction{}
}
