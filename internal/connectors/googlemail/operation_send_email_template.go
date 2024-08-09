package googlemail

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type sendTemplateMailOperationProps struct {
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

type SendTemplateMailOperation struct {
	options *sdk.OperationInfo
}

func NewSendTemplateMailOperation() sdk.IOperation {
	return &SendTemplateMailOperation{
		options: &sdk.OperationInfo{
			Name:        "Send Template Email",
			Description: "Send an personalized email.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
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
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *SendTemplateMailOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[sendTemplateMailOperationProps](ctx)
	gmailService, err := gmail.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	userProfile, err := gmailService.Users.GetProfile("me").Do()
	if err != nil {
		return nil, err
	}
	userEmail := userProfile.EmailAddress

	if !isValidEmail(input.To) {
		return nil, fmt.Errorf("invalid To email address: %s", input.To)
	}

	bccEmails := strings.Split(input.BCC, ",")
	for _, email := range bccEmails {
		email = strings.TrimSpace(email)
		if email != "" && !isValidEmail(email) {
			return nil, fmt.Errorf("invalid BCC email address: %s", email)
		}
	}

	ccEmails := strings.Split(input.CC, ",")
	for _, email := range ccEmails {
		email = strings.TrimSpace(email)
		if email != "" && !isValidEmail(email) {
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

func (c *SendTemplateMailOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *SendTemplateMailOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
