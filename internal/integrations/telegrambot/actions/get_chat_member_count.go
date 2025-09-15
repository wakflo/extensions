package actions

import (
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/telegrambot/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getChatMemberCountActionProps struct {
	ChatID string `json:"chat_id,omitempty"`
}

type GetChatMemberCountAction struct{}

func (a *GetChatMemberCountAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:          "get_chat_member_count",
		DisplayName: "Get Chat Member Count",
		Description: "Get the number of members in a chat.",
		Type:        core.ActionTypeAction,
		SampleOutput: map[string]any{
			"ok":     true,
			"result": 1543,
		},
		Settings: core.ActionSettings{},
	}
}

func (a *GetChatMemberCountAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_chat_administrators", "Get Chat Administrators")

	form.TextField("chat_id", "Chat ID").
		Placeholder("Enter chat ID").
		Required(true).
		HelpText("Unique identifier for the target chat of the target channel/group")

	schema := form.Build()

	return schema
}

func (a *GetChatMemberCountAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getChatMemberCountActionProps](ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %v", err)
	}

	if input.ChatID == "" {
		return nil, fmt.Errorf("chat_id is required")
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, fmt.Errorf("failed to get auth context: %v", err)
	}

	if authCtx.Extra == nil {
		return nil, fmt.Errorf("no authentication selected - please select a Telegram Bot connection")
	}

	token := authCtx.Extra["token"]
	if token == "" {
		return nil, fmt.Errorf("telegram bot token is empty - please check your connection configuration")
	}

	params := map[string]interface{}{
		"chat_id": input.ChatID,
	}

	response, err := shared.GetTelegramClient(token, "getChatMemberCount", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get member count for chat %s: %v", input.ChatID, err)
	}

	return response, nil
}

func (a *GetChatMemberCountAction) Auth() *core.AuthMetadata {
	return nil
}

func NewGetChatMemberCountAction() sdk.Action {
	return &GetChatAdministratorsAction{}
}
