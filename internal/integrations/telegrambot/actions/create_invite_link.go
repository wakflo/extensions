package actions

import (
	"fmt"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/telegrambot/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createInviteLinkActionProps struct {
	ChatID             string `json:"chat_id,omitempty"`
	Name               string `json:"name,omitempty"`
	ExpireDate         string `json:"expire_date,omitempty"`
	MemberLimit        int    `json:"member_limit,omitempty"`
	CreatesJoinRequest bool   `json:"creates_join_request,omitempty"`
}

type CreateInviteLinkAction struct{}

func (a *CreateInviteLinkAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_invite_link",
		DisplayName:   "Create Invite Link",
		Description:   "Create an invite link for a chat (group or channel). The bot must be an administrator in the chat.",
		Type:          core.ActionTypeAction,
		Documentation: createInviteLinkDocs,
		SampleOutput: map[string]any{
			"ok": true,
			"result": map[string]any{
				"invite_link": "https://t.me/+AbCdEfGhIjKlMnOp",
				"creator": map[string]any{
					"id":         123456789,
					"is_bot":     true,
					"first_name": "YourBot",
					"username":   "your_bot",
				},
				"creates_join_request": false,
				"is_primary":           false,
				"is_revoked":           false,
				"name":                 "Special Event Invite",
				"expire_date":          1614528967,
				"member_limit":         100,
			},
		},
		Settings: core.ActionSettings{},
	}
}

func (a *CreateInviteLinkAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_invite_link", "Create Invite Link")

	form.TextField("chat_id", "Chat ID").
		Placeholder("Enter chat ID or @username").
		Required(true).
		HelpText(`Unique identifier for the target chat or username of the target channel/supergroup.
		
**Note:** The bot must be an administrator in the chat with the following permissions:
- For groups: "can_invite_users" administrator right
- For channels: Administrator status`)

	form.TextField("name", "Invite Link Name").
		Placeholder("e.g., Special Event, New Members").
		Required(false).
		HelpText("Optional name for the invite link (0-32 characters)")

	form.TextField("expire_date", "Expiration Date").
		Placeholder("e.g., 2024-12-31T23:59:59Z or 1h30m").
		Required(false).
		HelpText(`When the link will expire. You can use:
- ISO 8601 format: 2024-12-31T23:59:59Z
- Duration from now: 1h, 30m, 24h, 7d
- Leave empty for no expiration`)

	form.NumberField("member_limit", "Member Limit").
		Placeholder("Maximum number of users").
		Required(false).
		DefaultValue(0).
		HelpText("Maximum number of users that can join via this link (1-99999). Use 0 for unlimited.")

	form.CheckboxField("creates_join_request", "Require Join Request").
		DefaultValue(false).
		HelpText("If true, users joining via this link will need to be approved by chat administrators")

	schema := form.Build()

	return schema
}

func (a *CreateInviteLinkAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createInviteLinkActionProps](ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %v", err)
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

	// Build parameters
	params := map[string]interface{}{
		"chat_id": input.ChatID,
	}

	// Add optional name
	if input.Name != "" {
		if len(input.Name) > 32 {
			return nil, fmt.Errorf("invite link name must be 32 characters or less")
		}
		params["name"] = input.Name
	}

	// Handle expiration date
	if input.ExpireDate != "" {
		expireTimestamp, err := parseExpireDate(input.ExpireDate)
		if err != nil {
			return nil, fmt.Errorf("invalid expire date format: %v", err)
		}
		params["expire_date"] = expireTimestamp
	}

	// Add member limit if specified
	if input.MemberLimit > 0 {
		if input.MemberLimit > 99999 {
			return nil, fmt.Errorf("member limit cannot exceed 99999")
		}
		params["member_limit"] = input.MemberLimit
	}

	// Add join request flag if true
	if input.CreatesJoinRequest {
		params["creates_join_request"] = true
		// Note: member_limit and creates_join_request are mutually exclusive
		if input.MemberLimit > 0 {
			return nil, fmt.Errorf("cannot use member_limit with creates_join_request")
		}
	}

	response, err := shared.GetTelegramClient(token, "createChatInviteLink", params)
	if err != nil {
		return nil, fmt.Errorf("failed to create invite link for chat %s: %v", input.ChatID, err)
	}

	return response, nil
}

func (a *CreateInviteLinkAction) Auth() *core.AuthMetadata {
	return nil
}

// Helper function to parse expire date
func parseExpireDate(expireStr string) (int64, error) {
	// First try to parse as duration (e.g., "1h", "30m", "7d")
	if duration, err := parseDuration(expireStr); err == nil {
		return time.Now().Add(duration).Unix(), nil
	}

	// Try to parse as ISO 8601 date
	if t, err := time.Parse(time.RFC3339, expireStr); err == nil {
		return t.Unix(), nil
	}

	// Try other common formats
	formats := []string{
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, expireStr); err == nil {
			return t.Unix(), nil
		}
	}

	return 0, fmt.Errorf("unrecognized date format. Use ISO 8601 (2024-12-31T23:59:59Z) or duration (1h, 30m, 7d)")
}

// Helper function to parse duration strings like "1h", "30m", "7d"
func parseDuration(s string) (time.Duration, error) {
	// Handle days specially since time.ParseDuration doesn't support them
	if len(s) > 1 && s[len(s)-1] == 'd' {
		days := 0
		if _, err := fmt.Sscanf(s, "%dd", &days); err == nil {
			return time.Duration(days) * 24 * time.Hour, nil
		}
	}

	// Try standard duration parsing for hours and minutes
	return time.ParseDuration(s)
}

func NewCreateInviteLinkAction() sdk.Action {
	return &CreateInviteLinkAction{}
}
