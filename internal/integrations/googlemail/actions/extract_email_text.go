package actions

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
	"golang.org/x/net/html"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type extractMailTextActionProps struct {
	MailID          string `json:"mail_id"`
	ExtractType     string `json:"extract_type"` // "text", "html", "both", "smart"
	IncludeHeaders  bool   `json:"include_headers"`
	IncludeMetadata bool   `json:"include_metadata"`
	ExtractLinks    bool   `json:"extract_links"`
	ExtractEmails   bool   `json:"extract_emails"`
	ExtractPhones   bool   `json:"extract_phones"`
	CleanHTML       bool   `json:"clean_html"`
}

type ExtractMailTextAction struct{}

func (a *ExtractMailTextAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "extract_mail_text",
		DisplayName:   "Extract Mail Text",
		Description:   "Extract and parse text content from Gmail emails, including plain text, HTML conversion, and data extraction.",
		Type:          core.ActionTypeAction,
		Documentation: extractMailTextDocs,
		SampleOutput: map[string]any{
			"success": true,
			"mail_id": "18d4a5b6c7e8f9g0",
			"content": map[string]any{
				"plain_text": "Hello,\n\nThis is the email content...",
				"subject":    "Meeting Tomorrow",
				"from":       "sender@example.com",
				"to":         []string{"recipient@example.com"},
				"date":       "Wed, 15 Jan 2024 10:30:00 +0000",
				"extracted": map[string]any{
					"links":  []string{"https://example.com"},
					"emails": []string{"contact@example.com"},
					"phones": []string{"+1-555-123-4567"},
				},
				"metadata": map[string]any{
					"word_count":       150,
					"char_count":       890,
					"has_attachments":  true,
					"attachment_count": 2,
				},
			},
		},
		Settings: core.ActionSettings{},
	}
}

func (a *ExtractMailTextAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("extract_mail_text", "Extract Mail Text")

	form.TextField("mail_id", "Mail ID").
		Placeholder("Enter message ID").
		Required(true).
		HelpText("The message ID of the email to extract text from")

	form.SelectField("extract_type", "Extract Type").
		Placeholder("Select extraction type").
		Required(false).
		AddOptions([]*smartform.Option{
			{Value: "smart", Label: "Smart (Best Available)"},
			{Value: "text", Label: "Plain Text Only"},
			{Value: "html", Label: "HTML Only"},
			{Value: "both", Label: "Both Text and HTML"},
		}...).
		DefaultValue("smart").
		HelpText("Type of content to extract from the email")

	form.CheckboxField("include_headers", "Include Headers").
		DefaultValue(true).
		HelpText("Include email headers (From, To, Subject, Date, etc.)")

	form.CheckboxField("include_metadata", "Include Metadata").
		DefaultValue(false).
		HelpText("Include metadata like word count, attachments info")

	form.CheckboxField("clean_html", "Clean HTML").
		DefaultValue(true).
		HelpText("Convert HTML to clean, readable plain text")

	form.CheckboxField("extract_links", "Extract Links").
		DefaultValue(false).
		HelpText("Extract all URLs from the email content")

	form.CheckboxField("extract_emails", "Extract Email Addresses").
		DefaultValue(false).
		HelpText("Extract all email addresses mentioned in the content")

	form.CheckboxField("extract_phones", "Extract Phone Numbers").
		DefaultValue(false).
		HelpText("Extract all phone numbers from the content")

	schema := form.Build()
	return schema
}

func (a *ExtractMailTextAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *ExtractMailTextAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[extractMailTextActionProps](ctx)
	if err != nil {
		return nil, err
	}

	if input.MailID == "" {
		return nil, errors.New("mail ID is required")
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	gmailService, err := gmail.NewService(context.Background(), option.WithTokenSource(*authCtx.TokenSource))
	if err != nil {
		return nil, err
	}

	message, err := gmailService.Users.Messages.Get("me", input.MailID).
		Format("full").
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch email: %v", err)
	}

	result := map[string]any{
		"success": true,
		"mail_id": input.MailID,
		"content": map[string]any{},
	}
	content := result["content"].(map[string]any)

	if input.IncludeHeaders {
		headers := extractGmailHeaders(message.Payload)
		content["subject"] = headers["Subject"]
		content["from"] = headers["From"]
		content["to"] = parseEmailList(headers["To"])
		content["cc"] = parseEmailList(headers["Cc"])
		content["date"] = headers["Date"]
		content["message_id"] = headers["Message-ID"]
		content["reply_to"] = headers["Reply-To"]
	}

	var plainText, htmlContent string

	switch input.ExtractType {
	case "text":
		plainText = extractGmailText(message.Payload, "text/plain")
		content["plain_text"] = plainText

	case "html":
		htmlContent = extractGmailText(message.Payload, "text/html")
		content["html_content"] = htmlContent

	case "both":
		plainText = extractGmailText(message.Payload, "text/plain")
		htmlContent = extractGmailText(message.Payload, "text/html")
		content["plain_text"] = plainText
		content["html_content"] = htmlContent

	default: // "smart"
		plainText = extractGmailText(message.Payload, "text/plain")
		if plainText == "" {
			htmlContent = extractGmailText(message.Payload, "text/html")
			if htmlContent != "" && input.CleanHTML {
				plainText = htmlToText(htmlContent)
			}
		}
		content["plain_text"] = plainText

		if input.CleanHTML && htmlContent != "" {
			content["cleaned_text"] = htmlToText(htmlContent)
		}
	}

	textForExtraction := plainText
	if textForExtraction == "" && htmlContent != "" {
		textForExtraction = htmlToText(htmlContent)
	}

	if input.ExtractLinks || input.ExtractEmails || input.ExtractPhones {
		extracted := map[string]any{}

		if input.ExtractLinks {
			extracted["links"] = extractLinks(textForExtraction)
		}

		if input.ExtractEmails {
			extracted["emails"] = extractEmailAddresses(textForExtraction)
		}

		if input.ExtractPhones {
			extracted["phones"] = extractPhoneNumbers(textForExtraction)
		}

		content["extracted"] = extracted
	}

	if input.IncludeMetadata {
		attachmentCount := 0
		for _, part := range message.Payload.Parts {
			if part.Filename != "" {
				attachmentCount++
			}
		}

		metadata := map[string]any{
			"word_count":       countWords(textForExtraction),
			"char_count":       len(textForExtraction),
			"has_attachments":  attachmentCount > 0,
			"attachment_count": attachmentCount,
			"thread_id":        message.ThreadId,
			"label_ids":        message.LabelIds,
			"snippet":          message.Snippet,
		}
		content["metadata"] = metadata
	}

	return result, nil
}

func extractGmailHeaders(payload *gmail.MessagePart) map[string]string {
	headers := make(map[string]string)
	for _, header := range payload.Headers {
		headers[header.Name] = header.Value
	}
	return headers
}

func extractGmailText(payload *gmail.MessagePart, mimeType string) string {
	// Check if this part matches the desired MIME type
	if payload.MimeType == mimeType && payload.Body != nil && payload.Body.Data != "" {
		data, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err != nil {
			return ""
		}
		return string(data)
	}

	// Recursively check parts
	for _, part := range payload.Parts {
		if text := extractGmailText(part, mimeType); text != "" {
			return text
		}
	}

	return ""
}

func parseEmailList(emailStr string) []string {
	if emailStr == "" {
		return []string{}
	}

	emails := strings.Split(emailStr, ",")
	for i, email := range emails {
		emails[i] = strings.TrimSpace(email)
	}

	return emails
}

func htmlToText(htmlContent string) string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return htmlContent
	}

	var textBuilder strings.Builder
	var extractText func(*html.Node)

	extractText = func(n *html.Node) {
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				textBuilder.WriteString(text)
				textBuilder.WriteString(" ")
			}
		}

		if n.Type == html.ElementNode {
			switch n.Data {
			case "br", "p", "div", "li":
				textBuilder.WriteString("\n")
			case "h1", "h2", "h3", "h4", "h5", "h6":
				textBuilder.WriteString("\n\n")
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}

	extractText(doc)

	text := textBuilder.String()
	text = regexp.MustCompile(`\n{3,}`).ReplaceAllString(text, "\n\n")
	text = regexp.MustCompile(` {2,}`).ReplaceAllString(text, " ")

	return strings.TrimSpace(text)
}

func extractLinks(text string) []string {
	urlRegex := regexp.MustCompile(`https?://[^\s<>"{}|\\^` + "`" + `\[\]]+`)
	matches := urlRegex.FindAllString(text, -1)

	seen := make(map[string]bool)
	var unique []string

	for _, match := range matches {
		match = strings.TrimRight(match, ".,;:!?")
		if !seen[match] {
			seen[match] = true
			unique = append(unique, match)
		}
	}

	return unique
}

func extractEmailAddresses(text string) []string {
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	matches := emailRegex.FindAllString(text, -1)

	seen := make(map[string]bool)
	var unique []string

	for _, match := range matches {
		lower := strings.ToLower(match)
		if !seen[lower] {
			seen[lower] = true
			unique = append(unique, match)
		}
	}

	return unique
}

func extractPhoneNumbers(text string) []string {
	phonePatterns := []*regexp.Regexp{
		regexp.MustCompile(`\+?[1-9]\d{0,2}[-.\s]?\(?\d{1,4}\)?[-.\s]?\d{1,4}[-.\s]?\d{1,9}`),
		regexp.MustCompile(`\(\d{3}\)\s*\d{3}[-.\s]?\d{4}`),
		regexp.MustCompile(`\d{3}[-.\s]?\d{3}[-.\s]?\d{4}`),
	}

	seen := make(map[string]bool)
	var phones []string

	for _, pattern := range phonePatterns {
		matches := pattern.FindAllString(text, -1)
		for _, match := range matches {
			cleaned := strings.TrimSpace(match)
			if len(cleaned) >= 10 && !seen[cleaned] {
				seen[cleaned] = true
				phones = append(phones, cleaned)
			}
		}
	}

	return phones
}

func countWords(text string) int {
	words := strings.Fields(text)
	return len(words)
}

func NewExtractMailTextAction() sdk.Action {
	return &ExtractMailTextAction{}
}
