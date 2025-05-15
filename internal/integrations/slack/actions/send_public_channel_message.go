package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/slack/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type sendPublicChannelMessageActionProps struct {
	Channel string `json:"channel"`
	Message string `json:"message"`
}

type SendPublicChannelMessageAction struct{}

func (a *SendPublicChannelMessageAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "send_public_channel_message",
		DisplayName:   "Send Public Channel Message",
		Description:   "Sends a message to a public channel in your chosen messaging platform, allowing you to share updates and information with team members or stakeholders.",
		Type:          core.ActionTypeAction,
		Documentation: sendPublicChannelMessageDocs,
		SampleOutput: map[string]interface{}{
			"name":       "slack-send-public-channel-message",
			"usage_mode": "operation",
			"message":    "Hello people in the public channel!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *SendPublicChannelMessageAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("send_public_channel_message", "Send Public Channel Message")

	// Define the function to get public channels
	getPublicChannels := func(ctx sdkcontext.DynamicFieldContext) (*core.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		client := shared.GetSlackClient(authCtx.Token.AccessToken)

		publicChannels, err := shared.GetChannels(client, "public_channel")
		if err != nil {
			return nil, err
		}

		options := make([]map[string]interface{}, 0, len(publicChannels))
		for _, channel := range publicChannels {
			options = append(options, map[string]interface{}{
				"value": channel.ID,
				"label": channel.Name,
			})
		}

		return ctx.Respond(options, len(options))
	}

	// Register the channel field with dynamic options
	form.SelectField("channel", "Public Channel").
		Placeholder("Select public channel where message will be sent").
		Required(true).
		HelpText("Select public channel where message will be sent").
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getPublicChannels)).
				RefreshOn("connection").
				GetDynamicSource(),
		)

	// Add message field
	shared.RegisterSharedLongMessageField(form)

	schema := form.Build()
	return schema
}

func (a *SendPublicChannelMessageAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[sendPublicChannelMessageActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	client := shared.GetSlackClient(authCtx.Token.AccessToken)

	message := input.Message
	channelID := input.Channel

	err = shared.SendMessage(client, message, channelID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"name":       "slack-send-public-channel-message",
		"usage_mode": "operation",
		"message":    message,
	}, nil
}

func (a *SendPublicChannelMessageAction) Auth() *core.AuthMetadata {
	return nil
}

func NewSendPublicChannelMessageAction() sdk.Action {
	return &SendPublicChannelMessageAction{}
}
