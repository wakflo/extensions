package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/discord/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type findChannelActionProps struct {
	GuildID     string `json:"guild-id"`
	ChannelName string `json:"channel-name"`
}

type FindChannelAction struct{}

// Metadata returns metadata about the action
func (a *FindChannelAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "find_channel",
		DisplayName:   "Find Channel",
		Description:   "Find a Discord channel by name and/or type in a server",
		Type:          core.ActionTypeAction,
		Documentation: findChannelDocs,
		SampleOutput: map[string]any{
			"id":                    "938177461736890368",
			"type":                  0,
			"name":                  "general",
			"guild_id":              "857347647235678912",
			"position":              1,
			"topic":                 "General discussion channel",
			"nsfw":                  false,
			"parent_id":             "123456789012345678",
			"permission_overwrites": []interface{}{},
		},
		Settings: core.ActionSettings{},
	}
}

func (a *FindChannelAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("find_channel", "Find Channel")

	form.TextField("guild-id", "Guild ID").
		Required(true).
		HelpText("The ID of the Discord server to search in")

	form.TextField("channel-name", "Channel Name").
		Required(false).
		HelpText("The name of the channel to find (leave empty to search by type only)")

	schema := form.Build()
	return schema
}

func (a *FindChannelAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[findChannelActionProps](ctx)
	if err != nil {
		return nil, err
	}

	endpoint := "/guilds/" + input.GuildID + "/channels"

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["token"] == "" {
		return nil, errors.New("missing discord bot token")
	}

	response, err := shared.GetDiscordClient(authCtx.Extra["token"], endpoint, "GET", nil)
	if err != nil {
		return nil, err
	}

	var matchingChannels []interface{}

	for _, channel := range response {
		channelMap, ok := channel.(map[string]interface{})
		if !ok {
			continue
		}
		matches := true

		if input.ChannelName != "" {
			if channelName, exists := channelMap["name"]; exists {
				if channelName != input.ChannelName {
					matches = false
				}
			} else {
				matches = false
			}
		}

		if matches {
			matchingChannels = append(matchingChannels, channel)
		}
	}

	if len(matchingChannels) == 0 {
		return map[string]interface{}{
			"found":    false,
			"message":  "No channels found matching the criteria",
			"channels": []interface{}{},
		}, nil
	}

	if len(matchingChannels) == 1 {
		result := matchingChannels[0].(map[string]interface{})
		result["found"] = true
		return result, nil
	}

	// Multiple channels found, return all
	return map[string]interface{}{
		"found":    true,
		"message":  "Multiple channels found",
		"count":    len(matchingChannels),
		"channels": matchingChannels,
	}, nil
}

func (a *FindChannelAction) Auth() *core.AuthMetadata {
	return nil
}

func NewFindChannelAction() sdk.Action {
	return &FindChannelAction{}
}
