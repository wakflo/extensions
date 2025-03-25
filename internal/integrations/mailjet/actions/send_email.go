package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/wakflo/extensions/internal/integrations/mailjet/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
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

func (a *SendEmailAction) Name() string {
	return "Send Email"
}

func (a *SendEmailAction) Description() string {
	return "Send transactional or marketing emails to your contacts with personalized content."
}

func (a *SendEmailAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *SendEmailAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &sendEmailDocs,
	}
}

func (a *SendEmailAction) Icon() *string {
	icon := "mdi:email-send"
	return &icon
}

func (a *SendEmailAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"from_email": autoform.NewShortTextField().
			SetDisplayName("From Email").
			SetDescription("The sender's email address").
			SetRequired(true).Build(),

		"from_name": autoform.NewShortTextField().
			SetDisplayName("From Name").
			SetDescription("The sender's name").
			SetRequired(true).Build(),

		"to_email": autoform.NewShortTextField().
			SetDisplayName("To Email").
			SetDescription("The recipient's email address").
			SetRequired(true).Build(),

		"to_name": autoform.NewShortTextField().
			SetDisplayName("To Name").
			SetDescription("The recipient's name").
			SetRequired(false).Build(),

		"subject": autoform.NewShortTextField().
			SetDisplayName("Subject").
			SetDescription("The email subject line").
			SetRequired(true).Build(),

		"text_part": autoform.NewLongTextField().
			SetDisplayName("Text Content").
			SetDescription("The plain text content of the email").
			SetRequired(false).Build(),

		"html_part": autoform.NewLongTextField().
			SetDisplayName("HTML Content").
			SetDescription("The HTML content of the email").
			SetRequired(false).Build(),

		"cc": autoform.NewLongTextField().
			SetDisplayName("CC").
			SetDescription("Carbon copy recipients").
			SetRequired(false).Build(),

		"bcc": autoform.NewLongTextField().
			SetDisplayName("BCC").
			SetDescription("Blind carbon copy recipients").
			SetRequired(false).Build(),

		"template_id": autoform.NewNumberField().
			SetDisplayName("Template ID").
			SetDescription("ID of the template to use (if any)").
			SetRequired(false).Build(),
	}
}

func (a *SendEmailAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[sendEmailActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client, err := shared.GetMailJetClient(ctx.Auth.Extra["api_key"], ctx.Auth.Extra["secret_key"])
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

func (a *SendEmailAction) Auth() *sdk.Auth {
	return nil
}

func (a *SendEmailAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"Messages": []map[string]any{
			{
				"Status":   "success",
				"CustomID": "",
				"To": []map[string]any{
					{
						"Email":       "recipient@example.com",
						"MessageID":   "12345678901234567",
						"MessageUUID": "1ab23cd4-5e67-8f90-12gh-34ij56kl7890",
					},
				},
				"Cc":  []map[string]any{},
				"Bcc": []map[string]any{},
			},
		},
	}
}

func (a *SendEmailAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewSendEmailAction() sdk.Action {
	return &SendEmailAction{}
}
