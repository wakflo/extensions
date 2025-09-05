package actions

import (
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/discord/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type removeGuildMemberActionProps struct {
	GuildID string `json:"guild-id"`
	UserID  string `json:"user-id"`
	Reason  string `json:"reason"`
}

type RemoveGuildMemberAction struct{}

// Metadata returns metadata about the action
func (a *RemoveGuildMemberAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "remove_guild_member",
		DisplayName:   "Remove Guild Member",
		Description:   "Remove (kick) a member from a Discord server",
		Type:          core.ActionTypeAction,
		Documentation: removeGuildMemberDocs,
		SampleOutput: map[string]any{
			"success":  true,
			"user_id":  "123456789012345678",
			"guild_id": "857347647235678912",
			"action":   "removed",
			"reason":   "Violation of server rules",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *RemoveGuildMemberAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("remove_guild_member", "Remove Guild Member")

	shared.RegisterGuildsInput(form, "Guilds", "List of guilds", true)

	form.TextField("user-id", "User ID").
		Required(true).
		HelpText("The ID of the user to remove from the server")

	schema := form.Build()
	return schema
}

func (a *RemoveGuildMemberAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[removeGuildMemberActionProps](ctx)
	if err != nil {
		return nil, err
	}

	endpoint := "/guilds/" + input.GuildID + "/members/" + input.UserID

	response, err := shared.GetDiscordClient(ctx.Auth().Key, endpoint, "DELETE", nil)
	if err != nil {
		return nil, err
	}

	fmt.Println(response...)

	return map[string]interface{}{
		"success":  true,
		"user_id":  input.UserID,
		"guild_id": input.GuildID,
		"action":   "removed",
		"reason":   input.Reason,
		"message":  "Member successfully removed from guild",
	}, nil
}

func (a *RemoveGuildMemberAction) Auth() *core.AuthMetadata {
	return nil
}

func NewRemoveGuildMemberAction() sdk.Action {
	return &RemoveGuildMemberAction{}
}
