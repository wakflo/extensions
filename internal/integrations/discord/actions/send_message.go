package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/discord/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type sendMessageActionProps struct {
	ChannelID string        `json:"channel-id"`
	Content   string        `json:"content"`
	Embed     *MessageEmbed `json:"embed"`
}

type MessageEmbed struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Color       int    `json:"color,omitempty"`
	URL         string `json:"url,omitempty"`
}

type SendMessageAction struct{}

// Metadata returns metadata about the action
func (a *SendMessageAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "send_message",
		DisplayName:   "Send Message",
		Description:   "Send a message to a specified Discord channel or user",
		Type:          core.ActionTypeAction,
		Documentation: sendMessageDocs,
		SampleOutput:  map[string]any{},
		Settings:      core.ActionSettings{},
	}
}

func (a *SendMessageAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("send_message", "Send Message")

	form.TextField("channel-id", "Channel ID").
		Required(true).
		HelpText("The ID of the Discord channel where the message will be sent")

	form.TextField("content", "Message Content").
		Required(true).
		HelpText("The content of the message to send")

	schema := form.Build()
	return schema
}

func (a *SendMessageAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[sendMessageActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Prepare the payload for Discord API
	payload := map[string]interface{}{
		"content": input.Content,
	}

	endpoint := "/channels/" + input.ChannelID + "/messages"

	response, err := shared.GetDiscordClient(ctx.Auth().Key, endpoint, "POST", payload)
	if err != nil {
		return nil, err
	}

	if len(response) > 0 {
		return response[0], nil
	}

	return map[string]interface{}{}, nil
}

func (a *SendMessageAction) Auth() *core.AuthMetadata {
	return nil
}

func NewSendMessageAction() sdk.Action {
	return &SendMessageAction{}
}
