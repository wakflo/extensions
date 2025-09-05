package triggers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/discord/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type messageReceivedTriggerProps struct {
	ChannelID string `json:"channel-id"`
	Content   string `json:"content,omitempty"`
}

type MessageReceivedTrigger struct{}

// Metadata returns metadata about the trigger
func (t *MessageReceivedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "message_received",
		DisplayName:   "Message Received",
		Description:   "Triggered when a new message is received in a Discord channel, with optional content filtering to monitor specific keywords or phrases.",
		Type:          core.TriggerTypePolling,
		Documentation: messageReceivedDocs,
		SampleOutput: map[string]any{
			"id":         "1095737357478123456",
			"type":       0,
			"content":    "Hello, this is a message!",
			"channel_id": "857347647235678912",
			"author": map[string]any{
				"id":            "123456789012345678",
				"username":      "User123",
				"discriminator": "1234",
				"avatar":        "abcdef123456",
			},
			"timestamp": "2023-04-12T12:34:56.789000+00:00",
		},
	}
}

// Props returns the schema for the trigger's input configuration
func (t *MessageReceivedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("message_received", "Message Received")

	shared.RegisterChannelsInput(form, "Channels", "The Discord channel to monitor", true)

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

func (t *MessageReceivedTrigger) Execute(ctx sdkcontext.ExecuteContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[messageReceivedTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	if input.ChannelID == "" {
		return nil, errors.New("Channel ID must be provided")
	}

	lastRunTime, err := ctx.GetMetadata("lastRun")
	if err != nil {
		return nil, err
	}

	var afterParam string
	if lastRunTime != nil {
		lastRunTimeValue := lastRunTime.(*time.Time)
		// Convert to snowflake timestamp for Discord
		// Discord uses snowflake IDs where the first 42 bits are milliseconds since Discord epoch (2015-01-01)
		discordEpoch := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
		millisSinceEpoch := lastRunTimeValue.Sub(discordEpoch).Milliseconds()
		// Shift by 22 bits to create a snowflake timestamp
		snowflake := millisSinceEpoch << 22
		afterParam = fmt.Sprintf("&after=%d", snowflake)
	} else {
		// If no last run time, just get the last 50 messages
		afterParam = ""
	}

	// Build the endpoint URL
	endpoint := fmt.Sprintf("/channels/%s/messages?limit=50%s", input.ChannelID, afterParam)

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	response, err := shared.GetDiscordClient(authCtx.Token.AccessToken, endpoint, "GET", nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching messages: %v", err)
	}

	messages := response

	var filteredMessages []interface{}
	for _, msg := range messages {
		message, ok := msg.(map[string]interface{})
		if !ok {
			continue
		}

		if channelID, ok := message["channel_id"].(string); ok {
			if channelID != input.ChannelID {
				continue
			}
		}

		filteredMessages = append(filteredMessages, message)
	}

	if len(filteredMessages) == 0 {
		return []interface{}{}, nil
	}

	return filteredMessages, nil
}

func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func (t *MessageReceivedTrigger) Criteria(ctx context.Context) core.TriggerCriteria {
	return core.TriggerCriteria{}
}

func NewMessageReceivedTrigger() sdk.Trigger {
	return &MessageReceivedTrigger{}
}
