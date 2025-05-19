package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/mailjet/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type sendEmailActionProps struct {
	FromEmail      string `json:"from_email"`
	FromName       string `json:"from_name"`
	ToEmail        string `json:"to_email"`
	ToName         string `json:"to_name"`
	Subject        string `json:"subject"`
	TextPart       string `json:"text_part"`
	HTMLPart       string `json:"html_part"`
	CC             string `json:"cc,omitempty"`
	BCC            string `json:"bcc,omitempty"`
	TemplateID     int    `json:"template_id,omitempty"`
	TemplateVars   string `json:"template_vars,omitempty"`
	TrackOpens     bool   `json:"track_opens"`
	TrackClicks    bool   `json:"track_clicks"`
	CustomID       string `json:"custom_id,omitempty"`
	CustomCampaign string `json:"custom_campaign,omitempty"`
}

type SendEmailAction struct{}

func (a *SendEmailAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "send_email",
		DisplayName:   "Send Email",
		Description:   "Send transactional or marketing emails to your contacts with personalized content.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: sendEmailDocs,
		SampleOutput: map[string]any{
			"MessageID": "123456",
			"Status":    "success",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *SendEmailAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("send_email", "Send Email")

	form.TextField("from_email", "From Email").
		Placeholder("Enter From Email").
		Required(true).
		HelpText("The sender's email address")

	form.TextField("from_name", "From Name").
		Placeholder("Enter From Name").
		Required(true).
		HelpText("The sender's name")

	form.TextField("to_email", "To Email").
		Placeholder("Enter To Email").
		Required(true).
		HelpText("The recipient's email address")

	form.TextField("to_name", "To Name").
		Placeholder("Enter To Name").
		Required(false).
		HelpText("The recipient's name")

	form.TextField("subject", "Subject").
		Placeholder("Enter Subject").
		Required(true).
		HelpText("The email subject line")

	form.TextareaField("text_part", "Text Content").
		Placeholder("Enter Text Content").
		Required(false).
		HelpText("The plain text content of the email")

	form.TextareaField("html_part", "HTML Content").
		Placeholder("Enter HTML Content").
		Required(false).
		HelpText("The HTML content of the email")

	form.TextField("cc", "CC").
		Placeholder("Enter CC").
		Required(false).
		HelpText("Carbon copy recipients")

	form.TextField("bcc", "BCC").
		Placeholder("Enter BCC").
		Required(false).
		HelpText("Blind carbon copy recipients")

	form.NumberField("template_id", "Template ID").
		Placeholder("Enter Template ID").
		Required(false).
		HelpText("ID of the template to use (if any)")

	schema := form.Build()

	return schema
}

func (a *SendEmailAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[sendEmailActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	client, err := shared.GetMailJetClient(authCtx.Extra["api_key"], authCtx.Extra["secret_key"])
	if err != nil {
		return nil, err
	}

	if input.TextPart == "" && input.HTMLPart == "" && input.TemplateID <= 0 {
		return nil, fmt.Errorf("at least one of Text Content, HTML Content, or Template ID must be provided")
	}

	message := map[string]interface{}{
		"From": map[string]interface{}{
			"Email": input.FromEmail,
			"Name":  input.FromName,
		},
		"To": []map[string]interface{}{
			{
				"Email": input.ToEmail,
				"Name":  input.ToName,
			},
		},
		"Subject": input.Subject,
	}

	// Add content parts if provided
	if input.TextPart != "" {
		message["TextPart"] = input.TextPart
	}

	if input.HTMLPart != "" {
		message["HTMLPart"] = input.HTMLPart
	}

	// Add CC recipients if specified - FIXED
	if input.CC != "" {
		// Split by commas if multiple emails are provided
		emails := strings.Split(input.CC, ",")
		cc := make([]map[string]interface{}, 0, len(emails))

		for _, email := range emails {
			email = strings.TrimSpace(email)
			if email != "" {
				cc = append(cc, map[string]interface{}{
					"Email": email,
				})
			}
		}

		if len(cc) > 0 {
			message["Cc"] = cc
		}
	}

	// Add BCC recipients if specified - FIXED
	if input.BCC != "" {
		// Split by commas if multiple emails are provided
		emails := strings.Split(input.BCC, ",")
		bcc := make([]map[string]interface{}, 0, len(emails))

		for _, email := range emails {
			email = strings.TrimSpace(email)
			if email != "" {
				bcc = append(bcc, map[string]interface{}{
					"Email": email,
				})
			}
		}

		if len(bcc) > 0 {
			message["Bcc"] = bcc
		}
	}

	// Add template if specified
	if input.TemplateID > 0 {
		message["TemplateID"] = input.TemplateID
		message["TemplateLanguage"] = true

		if input.TemplateVars != "" {
			var templateVars map[string]interface{}
			if err := json.Unmarshal([]byte(input.TemplateVars), &templateVars); err != nil {
				return nil, fmt.Errorf("invalid template variables JSON: %v", err)
			}
			message["Variables"] = templateVars
		}
	}

	payload := map[string]interface{}{
		"Messages": []map[string]interface{}{message},
	}

	var result map[string]interface{}
	err = client.Request(http.MethodPost, "/v3.1/send", payload, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *SendEmailAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewSendEmailAction() sdk.Action {
	return &SendEmailAction{}
}
