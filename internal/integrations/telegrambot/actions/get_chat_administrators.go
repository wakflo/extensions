package actions

import (
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/telegrambot/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getChatAdministratorsActionProps struct {
	ChatID string `json:"chat_id,omitempty"`
}

type GetChatAdministratorsAction struct{}

func (a *GetChatAdministratorsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_chat_administrators",
		DisplayName:   "Get Chat Administrators",
		Description:   "Get a list of administrators in a chat.",
		Type:          core.ActionTypeAction,
		Documentation: getChatAdministratorsDocs,
		SampleOutput: map[string]any{
			"ok": true,
			"result": []map[string]any{
				{
					"status": "creator",
					"user": map[string]any{
						"id":         123456789,
						"is_bot":     false,
						"first_name": "Alice",
						"username":   "alice",
					},
					"is_anonymous": false,
				},
				{
					"status": "administrator",
					"user": map[string]any{
						"id":         987654321,
						"is_bot":     false,
						"first_name": "Bob",
						"username":   "bob",
					},
					"can_be_edited":        true,
					"can_manage_chat":      true,
					"can_delete_messages":  true,
					"can_restrict_members": true,
					"can_promote_members":  false,
					"can_change_info":      true,
					"can_invite_users":     true,
					"can_pin_messages":     true,
					"is_anonymous":         false,
					"custom_title":         "Moderator",
				},
			},
		},
		Settings: core.ActionSettings{},
	}
}

func (a *GetChatAdministratorsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_chat_administrators", "Get Chat Administrators")

	form.TextField("chat_id", "Chat ID").
		Placeholder("Enter chat ID or @username").
		Required(true).
		HelpText("Unique identifier for the target chat or username of the target channel/group")

	schema := form.Build()

	return schema
}

func (a *GetChatAdministratorsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getChatAdministratorsActionProps](ctx)
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

	response, err := shared.GetTelegramClient(token, "getChatAdministrators", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat administrators for chat %s: %v", input.ChatID, err)
	}

	return response, nil
}

func (a *GetChatAdministratorsAction) Auth() *core.AuthMetadata {
	return nil
}

func NewGetChatAdministratorsAction() sdk.Action {
	return &GetChatAdministratorsAction{}
}
