package actions

import (
	"errors"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/telegrambot/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type sendPhotoActionProps struct {
	ChatID              string `json:"chat_id,omitempty"`
	PhotoURL            string `json:"photo_url,omitempty"`
	Caption             string `json:"caption,omitempty"`
	ParseMode           string `json:"parse_mode,omitempty"`
	DisableNotification bool   `json:"disable_notification,omitempty"`
	ReplyToMessageID    string `json:"reply_to_message_id,omitempty"`
}

type SendPhotoAction struct{}

func (a *SendPhotoAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "send_photo",
		DisplayName:   "Send Photo",
		Description:   "Send a photo with optional caption to a Telegram chat.",
		Type:          core.ActionTypeAction,
		Documentation: sendPhotoDocs,
		SampleOutput: map[string]any{
			"ok": true,
			"result": map[string]any{
				"message_id": 123,
				"from": map[string]any{
					"id":         123456789,
					"is_bot":     true,
					"first_name": "YourBot",
					"username":   "your_bot",
				},
				"chat": map[string]any{
					"id":         987654321,
					"first_name": "User",
					"username":   "user",
					"type":       "private",
				},
				"date": 1614528967,
				"photo": []map[string]any{
					{
						"file_id":        "AgACAgIAAxkBAAIBbF_JdZUAAT9c8D8vfNm0a3brFPTeAAJarzEbvG9gSPXTIo1lv9QvbFYQnqAuAAMBAAMCAANtAAPa7QACHgQ",
						"file_unique_id": "AQADbFYQnqAuAAMa7QAC",
						"file_size":      1428,
						"width":          90,
						"height":         90,
					},
					{
						"file_id":        "AgACAgIAAxkBAAIBbF_JdZUAAT9c8D8vfNm0a3brFPTeAAJarzEbvG9gSPXTIo1lv9QvbFYQnqAuAAMBAAMCAAN4AANd7gACHgQ",
						"file_unique_id": "AQADbFYQnqAuAANd7gAC",
						"file_size":      19493,
						"width":          320,
						"height":         320,
					},
				},
				"caption": "A beautiful image",
			},
		},
		Settings: core.ActionSettings{},
	}
}

func (a *SendPhotoAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("send_photo", "Send Photo")

	form.TextField("chat_id", "Chat ID").
		Placeholder("Enter chat ID").
		Required(true).
		HelpText("Unique identifier for the target chat or username of the target channel/group/user")

	form.TextField("photo_url", "Photo URL").
		Placeholder("https://example.com/photo.jpg").
		Required(true).
		HelpText("URL of the photo to send")

	form.TextareaField("caption", "Caption").
		Placeholder("Enter photo caption").
		Required(false).
		HelpText("Photo caption (may also be used when resending photos by file_id)")

	form.SelectField("parse_mode", "Parse Mode").
		Placeholder("Select parse mode").
		Required(false).
		AddOptions([]*smartform.Option{
			{Value: "markdownV2", Label: "Markdown"},
			{Value: "HTML", Label: "HTML"},
		}...).
		DefaultValue("markdownV2").
		HelpText("Mode for parsing entities in the caption")

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

func (a *SendPhotoAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[sendPhotoActionProps](ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %v", err)
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, fmt.Errorf("failed to get auth context: %v", err)
	}

	var token string
	if authCtx != nil && authCtx.Extra != nil {
		token = authCtx.Extra["token"]
	}

	if authCtx.Extra == nil {
		return nil, errors.New("no authentication selected - please select a Telegram Bot connection")
	}

	if token == "" {
		return nil, errors.New("telegram bot token is empty - please check your connection configuration")
	}

	params := map[string]interface{}{
		"chat_id": input.ChatID,
		"photo":   input.PhotoURL,
	}

	if input.Caption != "" {
		params["caption"] = input.Caption
	}

	if input.ParseMode != "" {
		params["parse_mode"] = input.ParseMode
	}

	if input.DisableNotification {
		params["disable_notification"] = input.DisableNotification
	}

	if input.ReplyToMessageID != "" {
		params["reply_to_message_id"] = input.ReplyToMessageID
	}

	response, err := shared.GetTelegramClient(token, "sendPhoto", params)
	if err != nil {
		return nil, fmt.Errorf("failed to send photo to chat %s: %v", input.ChatID, err)
	}

	return response, nil
}

func (a *SendPhotoAction) Auth() *core.AuthMetadata {
	return nil
}

func NewSendPhotoAction() sdk.Action {
	return &SendPhotoAction{}
}
