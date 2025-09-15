package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/telegrambot/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type sendMessageActionProps struct {
	ChatID                string `json:"chat_id,omitempty"`
	Text                  string `json:"text,omitempty"`
	ParseMode             string `json:"parse_mode,omitempty"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview,omitempty"`
	DisableNotification   bool   `json:"disable_notification,omitempty"`
	ReplyToMessageID      string `json:"reply_to_message_id,omitempty"`
}

type SendMessageAction struct{}

func (a *SendMessageAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "send_message",
		DisplayName:   "Send Message",
		Description:   "Send a text message to a specific Telegram chat, user, group, or channel.",
		Type:          core.ActionTypeAction,
		Documentation: sendMessageDocs,
		SampleOutput: map[string]any{
			"ok": true,
			"result": map[string]any{
				"message_id": "123",
				"from": map[string]any{
					"id":         "123456789",
					"is_bot":     true,
					"first_name": "YourBot",
					"username":   "your_bot",
				},
				"chat": map[string]any{
					"id":         "987654321",
					"first_name": "User",
					"username":   "user",
					"type":       "private",
				},
				"date": "1614528967",
				"text": "Hello from Wakflo!",
			},
		},
		Settings: core.ActionSettings{},
	}
}

func (a *SendMessageAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("send_message", "Send Message")

	form.TextField("chat_id", "Chat ID").
		Placeholder("Enter chat ID").
		Required(true).
		HelpText(`
1. Start a chat with @userinfobot or @getidsbot on Telegram.
2. Start a conversation with the bot.

**Note: Remember to initiate the chat with the bot, or you'll get an error for "chat not found.**`)

	form.TextareaField("text", "Message Text").
		Placeholder("Enter message text").
		Required(true).
		HelpText("Text of the message to be sent")

	form.SelectField("parse_mode", "Parse Mode").
		Placeholder("Select parse mode").
		Required(false).
		AddOptions([]*smartform.Option{
			{Value: "markdownV2", Label: "Markdown"},
			{Value: "HTML", Label: "HTML"},
		}...).
		DefaultValue("markdownV2").
		HelpText("Mode for parsing entities in the message text")

	form.CheckboxField("disable_web_page_preview", "Disable Web Page Preview").
		DefaultValue(false).
		HelpText("Disables link previews for links in this message")

	form.CheckboxField("disable_notification", "Disable Notification").
		DefaultValue(false).
		HelpText("Sends the message silently. Users will receive a notification with no sound.")

	form.TextField("reply_to_message_id", "Reply to Message ID").
		Placeholder("Enter message ID to reply to").
		Required(false).
		HelpText("If the message is a reply, ID of the original message")

	schema := form.Build()

	return schema
}

func (a *SendMessageAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[sendMessageActionProps](ctx)
	if err != nil {
		return nil, err
	}

	params := map[string]interface{}{
		"chat_id": input.ChatID,
		"text":    input.Text,
	}

	if input.ParseMode != "" {
		params["parse_mode"] = input.ParseMode
	}

	if input.DisableWebPagePreview {
		params["disable_web_page_preview"] = true
	}

	if input.DisableNotification {
		params["disable_notification"] = true
	}

	if input.ReplyToMessageID != "" {
		params["reply_to_message_id"] = input.ReplyToMessageID
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	response, err := shared.GetTelegramClient(authCtx.Extra["token"], "sendMessage", params)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *SendMessageAction) Auth() *core.AuthMetadata {
	return nil
}

func NewSendMessageAction() sdk.Action {
	return &SendMessageAction{}
}
