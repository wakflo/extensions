package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/discord/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createChannelActionProps struct {
	GuildID string `json:"guild-id"`
	Name    string `json:"name"`
}

type CreateChannelAction struct{}

// Metadata returns metadata about the action
func (a *CreateChannelAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_channel",
		DisplayName:   "Create Channel",
		Description:   "Create a new text or voice channel in a Discord server",
		Type:          core.ActionTypeAction,
		Documentation: createChannelDocs,
		SampleOutput: map[string]any{
			"id":       "938177461736890368",
			"type":     "0",
			"name":     "new-channel",
			"guild_id": "857347647235678912",
			"position": "5",
			"topic":    "Channel created by Wakflo automation",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *CreateChannelAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_channel", "Create Channel")

	shared.RegisterGuildsInput(form, "Guilds", "List of guilds", true)

	form.TextField("name", "Channel Name").
		Required(true).
		HelpText("The name of the new channel")

	schema := form.Build()
	return schema
}

func (a *CreateChannelAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createChannelActionProps](ctx)
	if err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"name": input.Name,
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["token"] == "" {
		return nil, errors.New("missing discord bot token")
	}

	endpoint := "/guilds/" + input.GuildID + "/channels"

	response, err := shared.GetDiscordClient(authCtx.Extra["token"], endpoint, "POST", payload)
	if err != nil {
		return nil, err
	}

	if len(response) > 0 {
		return response[0], nil
	}

	return map[string]interface{}{}, nil
}

func (a *CreateChannelAction) Auth() *core.AuthMetadata {
	return nil
}

func NewCreateChannelAction() sdk.Action {
	return &CreateChannelAction{}
}
