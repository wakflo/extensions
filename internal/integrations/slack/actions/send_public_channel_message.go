package actions

import (
	"github.com/wakflo/extensions/internal/integrations/slack/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type sendPublicChannelMessageActionProps struct {
	Channel string `json:"channel"`
	Message string `json:"message"`
}

type SendPublicChannelMessageAction struct{}

func (a *SendPublicChannelMessageAction) Name() string {
	return "Send Public Channel Message"
}

func (a *SendPublicChannelMessageAction) Description() string {
	return "Sends a message to a public channel in your chosen messaging platform, allowing you to share updates and information with team members or stakeholders."
}

func (a *SendPublicChannelMessageAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *SendPublicChannelMessageAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &sendPublicChannelMessageDocs,
	}
}

func (a *SendPublicChannelMessageAction) Icon() *string {
	return nil
}

func (a *SendPublicChannelMessageAction) Properties() map[string]*sdkcore.AutoFormSchema {
	getPublicChannels := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		client := shared.GetSlackClient(ctx.Auth.AccessToken)

		publicChannels, err := shared.GetChannels(client, "public_channel")
		if err != nil {
			return nil, err
		}

		return ctx.Respond(publicChannels, len(publicChannels))
	}

	return map[string]*sdkcore.AutoFormSchema{
		"channel": autoform.NewDynamicField(sdkcore.String).
			SetDisplayName("Public Channel").
			SetDescription("Select public channel where message will be sent").
			SetDynamicOptions(&getPublicChannels).
			SetDependsOn([]string{"connection"}).
			SetRequired(true).
			Build(),
		"message": shared.SharedLongMessageAutoform,
	}
}

func (a *SendPublicChannelMessageAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[sendPublicChannelMessageActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client := shared.GetSlackClient(ctx.Auth.AccessToken)

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

func (a *SendPublicChannelMessageAction) Auth() *sdk.Auth {
	return nil
}

func (a *SendPublicChannelMessageAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"name":       "slack-send-public-channel-message",
		"usage_mode": "operation",
		"message":    "Hello people in the public channel!",
	}
}

func (a *SendPublicChannelMessageAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewSendPublicChannelMessageAction() sdk.Action {
	return &SendPublicChannelMessageAction{}
}
