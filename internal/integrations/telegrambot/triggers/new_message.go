package triggers

import (
	"context"
	"errors"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/telegrambot/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type messageReceivedTriggerProps struct {
	FilterChatType string `json:"filter_chat_type,omitempty"`
	FilterChatID   string `json:"filter_chat_id,omitempty"`
}

type MessageReceivedTrigger struct{}

// Metadata returns metadata about the trigger
func (t *MessageReceivedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "message_received",
		DisplayName:   "Message Received",
		Description:   "Triggered when your Telegram bot receives a new message from a user, allowing you to automate responses and workflows based on incoming messages.",
		Type:          core.TriggerTypePolling,
		Documentation: messageReceivedDocs,
		SampleOutput: map[string]any{
			"messages": []map[string]any{
				{
					"update_id": 123456789,
					"message": map[string]any{
						"message_id": 123,
						"from": map[string]any{
							"id":            987654321,
							"is_bot":        false,
							"first_name":    "User",
							"username":      "user",
							"language_code": "en",
						},
						"chat": map[string]any{
							"id":         987654321,
							"first_name": "User",
							"username":   "user",
							"type":       "private",
						},
						"date": 1614528967,
						"text": "Hello bot!",
					},
				},
			},
			"lastUpdateID": 123456789,
		},
	}
}

// Props returns the schema for the trigger's input configuration
func (t *MessageReceivedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("message_received", "Message Received")

	form.SelectField("filter_chat_type", "Filter by Chat Type").
		Placeholder("Select chat type").
		Required(false).
		AddOptions([]*smartform.Option{
			{Value: "private", Label: "Private Messages"},
			{Value: "group", Label: "Group Messages"},
			{Value: "supergroup", Label: "Supergroup Messages"},
			{Value: "channel", Label: "Channel Messages"},
		}...).
		DefaultValue("").
		HelpText("Filter messages by chat type. Leave empty to receive messages from all chat types.")

	form.TextField("filter_chat_id", "Filter by Chat ID").
		Placeholder("Enter specific chat ID").
		Required(false).
		HelpText("Optionally filter messages from a specific chat. Leave empty to receive messages from all chats.")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the trigger
func (t *MessageReceivedTrigger) Auth() *core.AuthMetadata {
	return nil
}

// Start initializes the trigger
func (t *MessageReceivedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Stop shuts down the trigger
func (t *MessageReceivedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the main trigger logic
func (t *MessageReceivedTrigger) Execute(ctx sdkcontext.ExecuteContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[messageReceivedTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	fmt.Println("MessageReceivedTrigger Execute called")

	// Get the last update ID from metadata
	lastUpdateIDMeta, err := ctx.GetMetadata("lastUpdateID")
	if err != nil {
		// If no metadata exists, start from 0
		lastUpdateIDMeta = nil
	}

	var offset int64 = 0
	if lastUpdateIDMeta != nil {
		if lastUpdateID, ok := lastUpdateIDMeta.(float64); ok {
			offset = int64(lastUpdateID) + 1
		}
	}

	// Get the last run time for additional filtering
	lastRunTime, err := ctx.GetMetadata("lastRun")
	if err != nil {
		lastRunTime = nil
	}

	params := map[string]interface{}{
		"offset":  offset,
		"limit":   100,
		"timeout": 0,
	}

	// Get the token from the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	response, err := shared.GetTelegramClient(authCtx.Token.AccessToken, "getUpdates", params)
	if err != nil {
		return nil, fmt.Errorf("error fetching updates: %v", err)
	}

	result, ok := response["result"].([]interface{})
	if !ok {
		return nil, errors.New("invalid response format: result field is not an array")
	}

	// If there are no new messages, return empty array
	if len(result) == 0 {
		return map[string]interface{}{
			"messages":     []interface{}{},
			"lastUpdateID": offset - 1,
		}, nil
	}

	// Get the latest update ID for next polling
	lastUpdate := result[len(result)-1].(map[string]interface{})
	newLastUpdateID := int64(lastUpdate["update_id"].(float64))

	// Transform updates to just return messages with filtering
	messages := make([]interface{}, 0)
	minTimestamp := int64(0)
	if lastRunTime != nil {
		if lastRunTimeValue, ok := lastRunTime.(*int64); ok {
			minTimestamp = *lastRunTimeValue
		}
	}

	for _, update := range result {
		updateObj := update.(map[string]interface{})

		// Only process message updates
		if msg, ok := updateObj["message"].(map[string]interface{}); ok {
			// Apply chat type filter if specified
			if input.FilterChatType != "" {
				if chat, ok := msg["chat"].(map[string]interface{}); ok {
					if chatType, ok := chat["type"].(string); ok {
						if chatType != input.FilterChatType {
							continue
						}
					}
				}
			}

			// Apply chat ID filter if specified
			if input.FilterChatID != "" {
				if chat, ok := msg["chat"].(map[string]interface{}); ok {
					if chatID, ok := chat["id"].(float64); ok {
						if fmt.Sprintf("%.0f", chatID) != input.FilterChatID {
							continue
						}
					}
				}
			}

			// Check if the message is newer than our last run
			if date, ok := msg["date"].(float64); ok {
				if int64(date) > minTimestamp {
					// Add the whole update object to our results
					messages = append(messages, updateObj)
				}
			} else {
				// If no date, include it anyway
				messages = append(messages, updateObj)
			}
		}
	}

	// Store the last update ID in metadata for next run
	if err := ctx.SetMetadata("lastUpdateID", newLastUpdateID); err != nil {
		// Log error but don't fail the trigger
		fmt.Printf("Failed to set lastUpdateID metadata: %v\n", err)
	}

	// Return both the messages and the last update ID
	return map[string]interface{}{
		"messages":     messages,
		"lastUpdateID": newLastUpdateID,
	}, nil
}

// Criteria returns the criteria for the trigger
func (t *MessageReceivedTrigger) Criteria(ctx context.Context) core.TriggerCriteria {
	return core.TriggerCriteria{}
}

func NewMessageReceivedTrigger() sdk.Trigger {
	return &MessageReceivedTrigger{}
}
