package actions

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/wakflo/extensions/internal/integrations/googlemail/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type sendEmailTemplateActionProps struct {
	To              string `json:"to"`
	TemplateSubject string `json:"temp-subject"`
	Subject         string `json:"subject"`
	Body            string `json:"body"`
	CC              string `json:"cc"`
	BCC             string `json:"bcc"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	FromName        string `json:"from"`
}

type SendEmailTemplateAction struct{}

func (a *SendEmailTemplateAction) Name() string {
	return "Send Email Template"
}

func (a *SendEmailTemplateAction) Description() string {
	return "Sends an email to one or more recipients using a pre-defined template. The template can include placeholders for dynamic data, such as variables and conditional statements. This action allows you to automate the sending of personalized emails as part of your workflow process."
}

func (a *SendEmailTemplateAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *SendEmailTemplateAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &sendEmailTemplateDocs,
	}
}

func (a *SendEmailTemplateAction) Icon() *string {
	return nil
}

func (a *SendEmailTemplateAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"temp-subject": autoform.NewShortTextField().
			SetDisplayName(" Template title").
			SetDescription("The template title.").
			SetRequired(true).
			Build(),
		"subject": autoform.NewShortTextField().
			SetDisplayName(" Subject").
			SetDescription("The template subject you want to send.").
			SetRequired(true).
			Build(),
		"to": autoform.NewShortTextField().
			SetDisplayName(" To").
			SetDescription("The receiver of the message address").
			SetRequired(true).
			Build(),
		"from": autoform.NewShortTextField().
			SetDisplayName(" From").
			SetDescription("The sender of the message address").
			SetRequired(false).
			Build(),
		"firstName": autoform.NewShortTextField().
			SetDisplayName("First Name").
			SetDescription("First name").
			SetRequired(true).
			Build(),
		"lastName": autoform.NewShortTextField().
			SetDisplayName("Last Name").
			SetDescription("Last name").
			SetRequired(false).
			Build(),
	}
}

func (a *SendEmailTemplateAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[sendEmailTemplateActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	gmailService, err := gmail.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	userProfile, err := gmailService.Users.GetProfile("me").Do()
	if err != nil {
		return nil, err
	}
	userEmail := userProfile.EmailAddress

	if !shared.IsValidEmail(input.To) {
		return nil, fmt.Errorf("invalid To email address: %s", input.To)
	}

	bccEmails := strings.Split(input.BCC, ",")
	for _, email := range bccEmails {
		email = strings.TrimSpace(email)
		if email != "" && !shared.IsValidEmail(email) {
			return nil, fmt.Errorf("invalid BCC email address: %s", email)
		}
	}

	ccEmails := strings.Split(input.CC, ",")
	for _, email := range ccEmails {
		email = strings.TrimSpace(email)
		if email != "" && !shared.IsValidEmail(email) {
			return nil, fmt.Errorf("invalid CC email address: %s", email)
		}
	}

	// Search for the template message
	query := "subject:" + input.TemplateSubject
	searchResult, err := gmailService.Users.Messages.List("me").Q(query).Do()
	if err != nil {
		return nil, errors.New("error searching for template")
	}
	if len(searchResult.Messages) == 0 {
		return nil, errors.New("no template found with subject: " + input.TemplateSubject)
	}

	// Use the first matching message as the template
	templateMsg, err := gmailService.Users.Messages.Get("me", searchResult.Messages[0].Id).Do()
	if err != nil {
		return nil, errors.New("error fetching template message")
	}

	// Decode the email body
	var emailBody string
	for _, part := range templateMsg.Payload.Parts {
		if part.MimeType == "text/html" {
			data, _ := base64.URLEncoding.DecodeString(part.Body.Data)
			emailBody = string(data)
			break
		}
	}
	if emailBody == "" {
		return nil, errors.New("no HTML body found in template message")
	}

	emailBody = strings.ReplaceAll(emailBody, "{{FirstName}}", input.FirstName)
	emailBody = strings.ReplaceAll(emailBody, "{{LastName}}", input.LastName)

	var fromField string
	if input.FromName != "" {
		fromField = fmt.Sprintf("%s <%s>", input.FromName, userEmail)
	} else {
		fromField = userEmail
	}

	// Use the modified template as the email body
	message := &gmail.Message{
		Raw: base64.URLEncoding.EncodeToString([]byte(
			"From: " + fromField + "\r\n" +
				"To: " + input.To + "\r\n" +
				"Subject: " + input.TemplateSubject + "\r\n" +
				"Content-Type: text/html; charset=UTF-8\r\n\r\n" +
				emailBody)),
	}

	// Send the message
	_, err = gmailService.Users.Messages.Send("me", message).Do()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"status": "Message sent successfully!",
	}, nil
}

func (a *SendEmailTemplateAction) Auth() *sdk.Auth {
	return nil
}

func (a *SendEmailTemplateAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *SendEmailTemplateAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewSendEmailTemplateAction() sdk.Action {
	return &SendEmailTemplateAction{}
}
