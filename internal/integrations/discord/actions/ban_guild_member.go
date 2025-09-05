package actions

import (
	"errors"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/discord/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type banGuildMemberActionProps struct {
	GuildID string `json:"guild-id"`
	UserID  string `json:"user-id"`
	Reason  string `json:"reason"`
}

type BanGuildMemberAction struct{}

// Metadata returns metadata about the action
func (a *BanGuildMemberAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "ban_guild_member",
		DisplayName:   "Ban Guild Member",
		Description:   "Ban a member from a Discord server and optionally delete their recent messages",
		Type:          core.ActionTypeAction,
		Documentation: banGuildMemberDocs,
		SampleOutput: map[string]any{
			"success":                true,
			"user_id":                "123456789012345678",
			"guild_id":               "857347647235678912",
			"action":                 "banned",
			"reason":                 "Spam and harassment",
			"delete_message_seconds": 604800,
		},
		Settings: core.ActionSettings{},
	}
}

func (a *BanGuildMemberAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("ban_guild_member", "Ban Guild Member")

	shared.RegisterGuildsInput(form, "Guilds", "List of guilds", true)

	form.TextField("user-id", "User ID").
		Required(true).
		HelpText("The ID of the user to ban from the server")

	form.TextareaField("reason", "Reason").
		Required(false).
		HelpText("Optional reason for banning the member (will appear in audit log)")

	schema := form.Build()
	return schema
}

func (a *BanGuildMemberAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[banGuildMemberActionProps](ctx)
	if err != nil {
		return nil, err
	}

	endpoint := "/guilds/" + input.GuildID + "/bans/" + input.UserID

	payload := map[string]interface{}{}

	if input.Reason != "" {
		payload["reason"] = input.Reason
	}
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["token"] == "" {
		return nil, errors.New("missing discord bot token")
	}

	response, err := shared.GetDiscordClient(authCtx.Extra["token"], endpoint, "PUT", payload)
	if err != nil {
		return nil, err
	}

	fmt.Println(response...)

	return map[string]interface{}{
		"success":  true,
		"user_id":  input.UserID,
		"guild_id": input.GuildID,
		"action":   "banned",
		"reason":   input.Reason,
		"message":  "Member successfully banned from guild",
	}, nil
}

func (a *BanGuildMemberAction) Auth() *core.AuthMetadata {
	return nil
}

func NewBanGuildMemberAction() sdk.Action {
	return &BanGuildMemberAction{}
}
