package actions

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/googlemail/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type sendEmailTemplateActionProps struct {
	To              string `json:"to"`
	TemplateSubject string `json:"temp-subject"`
	Subject         string `json:"subject"`
	Body            string `json:"body"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	FromName        string `json:"from"`
}

type SendEmailTemplateAction struct{}

// Metadata returns metadata about the action
func (a *SendEmailTemplateAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "send_email_template",
		DisplayName:   "Send Email Template",
		Description:   "Sends an email to one or more recipients using a pre-defined template. The template can include placeholders for dynamic data, such as variables and conditional statements. This action allows you to automate the sending of personalized emails as part of your workflow process.",
		Type:          core.ActionTypeAction,
		Documentation: sendEmailTemplateDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *SendEmailTemplateAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("send_email_template", "Send Email Template")

	form.TextField("temp-subject", "temp-subject").
		Placeholder("Template title").
		HelpText("The template title.").
		Required(true)

	form.TextField("subject", "subject").
		Placeholder("Subject").
		HelpText("The template subject you want to send.").
		Required(true)

	form.TextField("to", "to").
		Placeholder("To").
		HelpText("The receiver of the message address").
		Required(true)

	form.TextField("from", "from").
		Placeholder("From").
		HelpText("The sender of the message address").
		Required(false)

	form.TextField("firstName", "firstName").
		Placeholder("First Name").
		HelpText("First name").
		Required(true)

	form.TextField("lastName", "lastName").
		Placeholder("Last Name").
		HelpText("Last name").
		Required(false)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *SendEmailTemplateAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *SendEmailTemplateAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[sendEmailTemplateActionProps](ctx)
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

	userProfile, err := gmailService.Users.GetProfile("me").Do()
	if err != nil {
		return nil, err
	}
	userEmail := userProfile.EmailAddress

	if !shared.IsValidEmail(input.To) {
		return nil, fmt.Errorf("invalid To email address: %s", input.To)
	}

	// bccEmails := strings.Split(input.BCC, ",")
	// for _, email := range bccEmails {
	// 	email = strings.TrimSpace(email)
	// 	if email != "" && !shared.IsValidEmail(email) {
	// 		return nil, fmt.Errorf("invalid BCC email address: %s", email)
	// 	}
	// }

	// ccEmails := strings.Split(input.CC, ",")
	// for _, email := range ccEmails {
	// 	email = strings.TrimSpace(email)
	// 	if email != "" && !shared.IsValidEmail(email) {
	// 		return nil, fmt.Errorf("invalid CC email address: %s", email)
	// 	}
	// }

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

func NewSendEmailTemplateAction() sdk.Action {
	return &SendEmailTemplateAction{}
}
