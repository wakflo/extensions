package actions

import (
	"github.com/wakflo/extensions/internal/integrations/slack/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type sendPrivateChannelMessageActionProps struct {
	Channel string `json:"channel"`
	Message string `json:"message"`
}

type SendPrivateChannelMessageAction struct{}

func (a *SendPrivateChannelMessageAction) Name() string {
	return "Send Private Channel Message"
}

func (a *SendPrivateChannelMessageAction) Description() string {
	return "Send Private Channel Message: Sends a private message to a specified user in a designated channel within your workflow automation platform, allowing you to discreetly communicate with team members or stakeholders without broadcasting to the entire channel."
}

func (a *SendPrivateChannelMessageAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *SendPrivateChannelMessageAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &sendPrivateChannelMessageDocs,
	}
}

func (a *SendPrivateChannelMessageAction) Icon() *string {
	return nil
}

func (a *SendPrivateChannelMessageAction) Properties() map[string]*sdkcore.AutoFormSchema {
	getPrivateChannels := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		client := shared.GetSlackClient(ctx.Auth.AccessToken)

		publicChannels, err := shared.GetChannels(client, "private_channel")
		if err != nil {
			return nil, err
		}

		return ctx.Respond(publicChannels, len(publicChannels))
	}

	return map[string]*sdkcore.AutoFormSchema{
		"channel": autoform.NewDynamicField(sdkcore.String).
			SetDisplayName("Private Channel").
			SetDescription("Select private channel where message will be sent").
			SetDynamicOptions(&getPrivateChannels).
			SetDependsOn([]string{"connection"}).
			SetRequired(true).
			Build(),
		"message": shared.SharedLongMessageAutoform,
	}
}

func (a *SendPrivateChannelMessageAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[sendPrivateChannelMessageActionProps](ctx.BaseContext)
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
		"name":       "slack-send-private-channel-message",
		"usage_mode": "operation",
	}, nil
}

func (a *SendPrivateChannelMessageAction) Auth() *sdk.Auth {
	return nil
}

func (a *SendPrivateChannelMessageAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"name":       "slack-send-private-channel-message",
		"usage_mode": "operation",
		"message":    "Hello people in the private channel!",
	}
}

func (a *SendPrivateChannelMessageAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewSendPrivateChannelMessageAction() sdk.Action {
	return &SendPrivateChannelMessageAction{}
}
