package triggers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/googlemail/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type TriggerMode string

const (
	ModeNewEmail        TriggerMode = "new_email"
	ModeSearchMatch     TriggerMode = "search_match"
	ModeNewAttachment   TriggerMode = "new_attachment"
	ModeNewLabeled      TriggerMode = "new_labeled"
	ModeNewStarred      TriggerMode = "new_starred"
	ModeNewConversation TriggerMode = "new_conversation"
)

type newEmailTriggerProps struct {
	Mode         TriggerMode `json:"mode"`
	SearchQuery  string      `json:"searchQuery"` // For search mode
	Label        string      `json:"label"`       // For labeled mode
	IncludeSpam  bool        `json:"includeSpam"`
	IncludeTrash bool        `json:"includeTrash"`
	MaxResults   int         `json:"maxResults"`
	// Legacy fields for backward compatibility
	Subject string `json:"subject"`
	From    string `json:"from"`
}

type NewEmailTrigger struct{}

// ✅ REQUIRED: Metadata method
func (t *NewEmailTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "new_email",
		DisplayName:   "Gmail Trigger",
		Description:   "Triggers on various Gmail events: new emails, attachments, labels, stars, and custom searches",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: newEmailDocs,
		SampleOutput: map[string]any{
			"id":             "12345abcde",
			"threadId":       "thread123",
			"subject":        "Important Message",
			"from":           "sender@example.com",
			"to":             "recipient@example.com",
			"date":           "Mon, 01 Jan 2025 10:00:00 -0700",
			"snippet":        "This is the beginning of the email...",
			"labels":         []string{"INBOX", "IMPORTANT"},
			"hasAttachments": true,
			"attachments": []map[string]any{
				{
					"filename": "document.pdf",
					"mimeType": "application/pdf",
					"size":     102400,
				},
			},
			"isStarred":   false,
			"isImportant": true,
			"isUnread":    true,
		},
	}
}

func (t *NewEmailTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func (t *NewEmailTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *NewEmailTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("google-mail-trigger", "Gmail Trigger Configuration")

	form.SelectField("mode", "Trigger Mode").
		Required(true).
		AddOption("new_email", "New Email").
		AddOption("search_match", "New Email Matching Search").
		AddOption("new_attachment", "New Attachment").
		AddOption("new_labeled", "New Labeled Email").
		AddOption("new_starred", "New Starred Email").
		AddOption("new_conversation", "New Conversation").
		DefaultValue("new_email").
		HelpText("Select what type of Gmail event should trigger this workflow")

	form.TextField("searchQuery", "Search Query").
		Placeholder("is:unread has:attachment from:important@example.com").
		HelpText("Gmail search query syntax (same as Gmail search box)").
		VisibleWhenEquals("mode", "search_match").
		Required(false)

	form.TextField("label", "Label").
		Placeholder("INBOX, IMPORTANT, or custom label").
		HelpText("Gmail label to monitor").
		VisibleWhenEquals("mode", "new_labeled").
		Required(false)

	form.TextField("subject", "Subject Filter").
		Placeholder("Invoice").
		HelpText("Filter emails by subject (optional)").
		Required(false)

	form.TextField("from", "From Filter").
		Placeholder("sender@example.com").
		HelpText("Filter emails by sender (optional)").
		Required(false)

	form.CheckboxField("includeSpam", "Include Spam").
		DefaultValue(false).
		HelpText("Include messages from SPAM folder")

	form.CheckboxField("includeTrash", "Include Trash").
		DefaultValue(false).
		HelpText("Include messages from TRASH folder")

	form.NumberField("maxResults", "Max Results").
		Placeholder("50").
		HelpText("Maximum number of emails to return per check").
		DefaultValue(50).
		Required(false)

	schema := form.Build()
	return schema
}

func (t *NewEmailTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *NewEmailTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *NewEmailTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[newEmailTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	gmailService, err := gmail.NewService(context.Background(), option.WithTokenSource(*authCtx.TokenSource))
	if err != nil {
		return nil, err
	}

	// Get last run time
	var lastRunTime time.Time
	lastRun, err := ctx.GetMetadata("lastRun")
	if err == nil && lastRun != nil {
		lastRunTime = *lastRun.(*time.Time)
	} else {
		// First run - set to 24 hours ago
		lastRunTime = time.Now().Add(-24 * time.Hour)
	}

	// Build query based on mode
	query := t.buildQuery(input, lastRunTime)

	// Configure list call
	listCall := gmailService.Users.Messages.List("me").Q(query)

	if input.MaxResults > 0 {
		listCall.MaxResults(int64(input.MaxResults))
	} else {
		listCall.MaxResults(50)
	}

	// Include spam/trash if specified
	if input.IncludeSpam || input.IncludeTrash {
		listCall.IncludeSpamTrash(true)
	}

	messages, err := listCall.Do()
	if err != nil {
		return nil, err
	}

	// Process messages based on mode
	results, err := t.processMessages(gmailService, messages, input, lastRunTime)
	if err != nil {
		return nil, err
	}

	// Update lastRun metadata
	now := time.Now()
	ctx.SetMetadata("lastRun", &now)

	return results, nil
}

func (t *NewEmailTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *NewEmailTrigger) buildQuery(input *newEmailTriggerProps, lastRunTime time.Time) string {
	var queryParts []string

	queryParts = append(queryParts, fmt.Sprintf("after:%d", lastRunTime.Unix()))

	switch input.Mode {
	case ModeSearchMatch:
		if input.SearchQuery != "" {
			queryParts = append(queryParts, input.SearchQuery)
		}

	case ModeNewAttachment:
		queryParts = append(queryParts, "has:attachment")

	case ModeNewLabeled:
		if input.Label != "" {
			queryParts = append(queryParts, fmt.Sprintf("label:%s", input.Label))
		}

	case ModeNewStarred:
		queryParts = append(queryParts, "is:starred")

	case ModeNewConversation:
		queryParts = append(queryParts, "in:inbox")

	default:
		queryParts = append(queryParts, "in:inbox")
	}

	if input.Subject != "" {
		queryParts = append(queryParts, fmt.Sprintf("subject:%s", input.Subject))
	}
	if input.From != "" {
		queryParts = append(queryParts, fmt.Sprintf("from:%s", input.From))
	}

	return strings.Join(queryParts, " ")
}

func (t *NewEmailTrigger) processMessages(
	service *gmail.Service,
	messages *gmail.ListMessagesResponse,
	input *newEmailTriggerProps,
	lastRunTime time.Time,
) ([]map[string]interface{}, error) {
	results := make([]map[string]interface{}, 0)
	processedThreads := make(map[string]bool)

	for _, msg := range messages.Messages {
		email, err := service.Users.Messages.Get("me", msg.Id).Do()
		if err != nil {
			continue
		}

		if input.Mode == ModeNewConversation {
			if processedThreads[email.ThreadId] {
				continue
			}
			processedThreads[email.ThreadId] = true
		}

		if input.Mode == ModeNewStarred {
			starredTime := t.getStarredTime(email)
			if starredTime.Before(lastRunTime) {
				continue
			}
		}

		emailData := t.extractEmailData(email)

		if input.Mode == ModeNewAttachment {
			if !emailData["hasAttachments"].(bool) {
				continue
			}
		}

		results = append(results, emailData)
	}

	return results, nil
}

func (t *NewEmailTrigger) extractEmailData(email *gmail.Message) map[string]interface{} {
	headers := email.Payload.Headers

	attachments := t.extractAttachments(email.Payload)

	return map[string]interface{}{
		"id":             email.Id,
		"threadId":       email.ThreadId,
		"subject":        shared.GetHeader(headers, "Subject"),
		"from":           shared.GetHeader(headers, "From"),
		"to":             shared.GetHeader(headers, "To"),
		"cc":             shared.GetHeader(headers, "Cc"),
		"date":           shared.GetHeader(headers, "Date"),
		"snippet":        email.Snippet,
		"labels":         email.LabelIds,
		"hasAttachments": len(attachments) > 0,
		"attachments":    attachments,
		"isStarred":      t.hasLabel(email.LabelIds, "STARRED"),
		"isImportant":    t.hasLabel(email.LabelIds, "IMPORTANT"),
		"isUnread":       t.hasLabel(email.LabelIds, "UNREAD"),
		"sizeEstimate":   email.SizeEstimate,
	}
}

func (t *NewEmailTrigger) extractAttachments(payload *gmail.MessagePart) []map[string]interface{} {
	attachments := []map[string]interface{}{}

	var findAttachments func(*gmail.MessagePart)
	findAttachments = func(part *gmail.MessagePart) {
		if part.Filename != "" && part.Body.AttachmentId != "" {
			attachments = append(attachments, map[string]interface{}{
				"attachmentId": part.Body.AttachmentId,
				"filename":     part.Filename,
				"mimeType":     part.MimeType,
				"size":         part.Body.Size,
			})
		}

		for _, subPart := range part.Parts {
			findAttachments(subPart)
		}
	}

	findAttachments(payload)
	return attachments
}

func (t *NewEmailTrigger) hasLabel(labels []string, label string) bool {
	for _, l := range labels {
		if l == label {
			return true
		}
	}
	return false
}

func (t *NewEmailTrigger) getStarredTime(email *gmail.Message) time.Time {
	return time.Unix(email.InternalDate/1000, 0)
}

func NewNewEmailTrigger() sdk.Trigger {
	return &NewEmailTrigger{}
}
