package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/slack/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type sendDirectMessageActionProps struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

type SendDirectMessageAction struct{}

func (a *SendDirectMessageAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "send_direct_message",
		DisplayName:   "Send Direct Message",
		Description:   "Send Direct Message: Automatically sends a direct message to a specific Slack channel or user, allowing you to notify team members of important updates or trigger custom workflows.",
		Type:          core.ActionTypeAction,
		Documentation: sendDirectMessageDocs,
		SampleOutput: map[string]interface{}{
			"name":       "slack-send-direct-message",
			"usage_mode": "operation",
			"message":    "Hello",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *SendDirectMessageAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("send_direct_message", "Send Direct Message")

	getSlackUsers := func(ctx sdkcontext.DynamicFieldContext) (*core.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		client := shared.GetSlackClient(authCtx.Token.AccessToken)

		users, err := shared.GetUsers(client)
		if err != nil {
			return nil, err
		}

		options := make([]map[string]interface{}, 0, len(users))
		for _, user := range users {
			options = append(options, map[string]interface{}{
				"id":   user.ID,
				"name": user.RealName,
			})
		}

		return ctx.Respond(options, len(options))
	}

	form.SelectField("user", "User").
		Placeholder("Select user to send message").
		Required(true).
		HelpText("Select user to send message").
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getSlackUsers)).
				RefreshOn("connection").
				GetDynamicSource(),
		)

	shared.RegisterSharedLongMessageField(form)

	schema := form.Build()
	return schema
}

func (a *SendDirectMessageAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[sendDirectMessageActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	client := shared.GetSlackClient(authCtx.Token.AccessToken)

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

func (a *SendDirectMessageAction) Auth() *core.AuthMetadata {
	return nil
}

func NewSendDirectMessageAction() sdk.Action {
	return &SendDirectMessageAction{}
}
