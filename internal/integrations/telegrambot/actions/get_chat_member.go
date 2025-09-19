package actions

import (
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/telegrambot/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getChatMemberActionProps struct {
	ChatID string `json:"chat_id,omitempty"`
	UserID string `json:"user_id,omitempty"`
}

type GetChatMemberAction struct{}

func (a *GetChatMemberAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_chat_member",
		DisplayName:   "Get Chat Member",
		Description:   "Get information about a member of a chat including their status (creator, administrator, member, restricted, left, or kicked).",
		Type:          core.ActionTypeAction,
		Documentation: getChatMemberDocs,
		SampleOutput: map[string]any{
			"ok": true,
			"result": map[string]any{
				"status": "administrator",
				"user": map[string]any{
					"id":         987654321,
					"is_bot":     false,
					"first_name": "John",
					"last_name":  "Doe",
					"username":   "johndoe",
				},
				"can_be_edited":          true,
				"can_manage_chat":        true,
				"can_change_info":        true,
				"can_delete_messages":    true,
				"can_invite_users":       true,
				"can_restrict_members":   true,
				"can_pin_messages":       true,
				"can_promote_members":    false,
				"can_manage_video_chats": true,
				"is_anonymous":           false,
				"can_manage_voice_chats": true,
				"custom_title":           "Community Manager",
			},
		},
		Settings: core.ActionSettings{},
	}
}

func (a *GetChatMemberAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_chat_member", "Get Chat Member")

	form.TextField("chat_id", "Chat ID").
		Placeholder("Enter chat ID").
		Required(true).
		HelpText(`Unique identifier for the target chat or username of the target channel/group.
		
**Examples:**
- Group/Channel ID: -1001234567890
- Public username: @channelname`)

	form.TextField("user_id", "User ID").
		Placeholder("Enter user ID").
		Required(true).
		HelpText(`Unique identifier of the target user.
		
**How to get a User ID:**
1. Have the user message your bot
2. Use the "Get Updates" action to see their user ID
3. Or forward a message from the user to @userinfobot to get their ID`)

	schema := form.Build()

	return schema
}

func (a *GetChatMemberAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getChatMemberActionProps](ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %v", err)
	}

	// Validate inputs
	if input.ChatID == "" {
		return nil, fmt.Errorf("chat_id is required")
	}
	if input.UserID == "" {
		return nil, fmt.Errorf("user_id is required")
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
		"user_id": input.UserID,
	}

	response, err := shared.GetTelegramClient(token, "getChatMember", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat member info for user %s in chat %s: %v",
			input.UserID, input.ChatID, err)
	}

	return response, nil
}

func (a *GetChatMemberAction) Auth() *core.AuthMetadata {
	return nil
}

func NewGetChatMemberAction() sdk.Action {
	return &GetChatMemberAction{}
}
