package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/slack/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type sendPrivateChannelMessageActionProps struct {
	Channel string `json:"channel"`
	Message string `json:"message"`
}

type SendPrivateChannelMessageAction struct{}

func (a *SendPrivateChannelMessageAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "send_private_channel_message",
		DisplayName:   "Send Private Channel Message",
		Description:   "Send Private Channel Message: Sends a private message to a specified user in a designated channel within your workflow automation platform, allowing you to discreetly communicate with team members or stakeholders without broadcasting to the entire channel.",
		Type:          core.ActionTypeAction,
		Documentation: sendPrivateChannelMessageDocs,
		SampleOutput: map[string]interface{}{
			"name":       "slack-send-private-channel-message",
			"usage_mode": "operation",
			"message":    "Hello people in the private channel!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *SendPrivateChannelMessageAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("send_private_channel_message", "Send Private Channel Message")

	// Define the function to get private channels
	getPrivateChannels := func(ctx sdkcontext.DynamicFieldContext) (*core.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		client := shared.GetSlackClient(authCtx.Token.AccessToken)

		privateChannels, err := shared.GetChannels(client, "private_channel")
		if err != nil {
			return nil, err
		}

		// Convert SlackChannel array to proper options format
		options := make([]map[string]interface{}, 0, len(privateChannels))
		for _, channel := range privateChannels {
			options = append(options, map[string]interface{}{
				"value": channel.ID,
				"label": channel.Name,
			})
		}

		return ctx.Respond(options, len(options))
	}

	form.SelectField("channel", "Private Channel").
		Placeholder("Select private channel where message will be sent").
		Required(true).
		HelpText("Select private channel where message will be sent").
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getPrivateChannels)).
				RefreshOn("connection").
				GetDynamicSource(),
		)

	shared.RegisterSharedLongMessageField(form)

	schema := form.Build()
	return schema
}

func (a *SendPrivateChannelMessageAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[sendPrivateChannelMessageActionProps](ctx)
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
		"name":       "slack-send-private-channel-message",
		"usage_mode": "operation",
		"message":    message,
	}, nil
}

func (a *SendPrivateChannelMessageAction) Auth() *core.AuthMetadata {
	return nil
}

func NewSendPrivateChannelMessageAction() sdk.Action {
	return &SendPrivateChannelMessageAction{}
}
