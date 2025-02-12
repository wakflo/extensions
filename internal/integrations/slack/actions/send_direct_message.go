package actions

import (
	"github.com/wakflo/extensions/internal/integrations/slack/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type sendDirectMessageActionProps struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

type SendDirectMessageAction struct{}

func (a *SendDirectMessageAction) Name() string {
	return "Send Direct Message"
}

func (a *SendDirectMessageAction) Description() string {
	return "Send Direct Message: Automatically sends a direct message to a specific Slack channel or user, allowing you to notify team members of important updates or trigger custom workflows."
}

func (a *SendDirectMessageAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *SendDirectMessageAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &sendDirectMessageDocs,
	}
}

func (a *SendDirectMessageAction) Icon() *string {
	return nil
}

func (a *SendDirectMessageAction) Properties() map[string]*sdkcore.AutoFormSchema {
	getSlackUsers := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		client := shared.GetSlackClient(ctx.Auth.AccessToken)

		users, err := shared.GetUsers(client)
		if err != nil {
			return nil, err
		}

		return ctx.Respond(users, len(users))
	}

	return map[string]*sdkcore.AutoFormSchema{
		"user": autoform.NewDynamicField(sdkcore.String).
			SetDisplayName("User").
			SetDescription("Select user to send message").
			SetDynamicOptions(&getSlackUsers).
			SetDependsOn([]string{"connection"}).
			SetRequired(true).
			Build(),
		"message": shared.SharedLongMessageAutoform,
	}
}

func (a *SendDirectMessageAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[sendDirectMessageActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client := shared.GetSlackClient(ctx.Auth.AccessToken)

	message := input.Message
	userID := input.User
	err = shared.SendMessage(client, message, userID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"name":       "slack-send-direct-message",
		"usage_mode": "operation",
		"message":    message,
	}, nil
}

func (a *SendDirectMessageAction) Auth() *sdk.Auth {
	return nil
}

func (a *SendDirectMessageAction) SampleData() sdkcore.JSON {
	return map[string]interface{}{
		"name":       "slack-send-direct-message",
		"usage_mode": "operation",
		"message":    "Hello",
	}
}

func (a *SendDirectMessageAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewSendDirectMessageAction() sdk.Action {
	return &SendDirectMessageAction{}
}
