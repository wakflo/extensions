package actions

import (
	"context"
	"encoding/base64"
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

type sendEmailActionProps struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	CC      string `json:"cc"`
	BCC     string `json:"bcc"`
}

type SendEmailAction struct{}

func (a *SendEmailAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "send_email",
		DisplayName:   "Send Email",
		Description:   "Sends an email to one or more recipients using a customizable template and attachments.",
		Type:          core.ActionTypeAction,
		Documentation: sendEmailDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *SendEmailAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("send_email", "Send Email")

	form.TextField("subject", "subject").
		Placeholder("Subject").
		HelpText("The subject of the message").
		Required(true)

	form.TextField("to", "to").
		Placeholder("To").
		HelpText("The receiver of the message address").
		Required(true)

	form.TextareaField("body", "body").
		Placeholder("Body").
		HelpText("Email body").
		Required(true)

	form.TextField("bcc", "bcc").
		Placeholder("BCC").
		HelpText("Blind carbon copy recipients (comma-separated)").
		Required(false)

	form.TextField("cc", "cc").
		Placeholder("CC").
		HelpText("Carbon copy recipients (comma-separated)").
		Required(false)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *SendEmailAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *SendEmailAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[sendEmailActionProps](ctx)
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

	fromEmail := userProfile.EmailAddress

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

	message := &gmail.Message{
		Raw: base64.URLEncoding.EncodeToString([]byte(
			"From: " + fromEmail + "\r\n" +
				"To: " + input.To + "\r\n" +
				"Subject: " + input.Subject + "\r\n" +
				"Bcc: " + strings.Join(bccEmails, ",") + "\r\n" +
				"Cc: " + strings.Join(ccEmails, ",") + "\r\n\r\n" +
				input.Body)),
	}

	_, err = gmailService.Users.Messages.Send("me", message).Do()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Message": "Message Sent Successfully",
	}, nil
}

func NewSendEmailAction() sdk.Action {
	return &SendEmailAction{}
}
