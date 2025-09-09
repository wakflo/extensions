package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/discord/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type deleteChannelActionProps struct {
	ChannelID string `json:"channel-id"`
}

type DeleteChannelAction struct{}

// Metadata returns metadata about the action
func (a *DeleteChannelAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "delete_channel",
		DisplayName:   "Delete Channel",
		Description:   "Delete a Discord channel from a server",
		Type:          core.ActionTypeAction,
		Documentation: deleteChannelDocs,
		SampleOutput: map[string]any{
			"id":       "938177461736890368",
			"type":     "0",
			"name":     "deleted-channel",
			"guild_id": "857347647235678912",
			"deleted":  true,
		},
		Settings: core.ActionSettings{},
	}
}

func (a *DeleteChannelAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("delete_channel", "Delete Channel")

	shared.RegisterChannelsInput(form, "Channels", "List of channels", true)

	schema := form.Build()
	return schema
}

func (a *DeleteChannelAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[deleteChannelActionProps](ctx)
	if err != nil {
		return nil, err
	}

	endpoint := "/channels/" + input.ChannelID

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["token"] == "" {
		return nil, errors.New("missing discord bot token")
	}

	response, err := shared.GetDiscordClient(authCtx.Extra["token"], endpoint, "DELETE", nil)
	if err != nil {
		return nil, err
	}

	if len(response) > 0 {
		return response[0], nil
	}

	return map[string]interface{}{}, nil
}

func (a *DeleteChannelAction) Auth() *core.AuthMetadata {
	return nil
}

func NewDeleteChannelAction() sdk.Action {
	return &DeleteChannelAction{}
}
