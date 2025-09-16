package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/discord/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type renameChannelActionProps struct {
	ChannelID string `json:"channel-id"`
	NewName   string `json:"name"`
}

type RenameChannelAction struct{}

// Metadata returns metadata about the action
func (a *RenameChannelAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "rename_channel",
		DisplayName:   "Rename Channel",
		Description:   "Rename a Discord channel and optionally update its topic",
		Type:          core.ActionTypeAction,
		Documentation: renameChannelDocs,
		SampleOutput: map[string]any{
			"id":       "938177461736890368",
			"type":     "0",
			"name":     "new-channel-name",
			"guild_id": "857347647235678912",
			"topic":    "Updated channel topic",
			"position": "5",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *RenameChannelAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("rename_channel", "Rename Channel")

	shared.RegisterChannelsInput(form, "Channels", "List of channels", true)

	form.TextField("name", "New Channel Name").
		Required(true).
		HelpText("The new name for the channel (must be 1-100 characters)")

	schema := form.Build()
	return schema
}

func (a *RenameChannelAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[renameChannelActionProps](ctx)
	if err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"name": input.NewName,
	}

	endpoint := "/channels/" + input.ChannelID

	response, err := shared.GetDiscordClient(ctx.Auth().Key, endpoint, "PATCH", payload)
	if err != nil {
		return nil, err
	}

	if len(response) > 0 {
		return response[0], nil
	}

	return map[string]interface{}{}, nil
}

func (a *RenameChannelAction) Auth() *core.AuthMetadata {
	return nil
}

func NewRenameChannelAction() sdk.Action {
	return &RenameChannelAction{}
}
