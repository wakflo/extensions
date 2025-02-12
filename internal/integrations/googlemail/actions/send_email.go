package actions

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/wakflo/extensions/internal/integrations/googlemail/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
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

func (a *SendEmailAction) Name() string {
	return "Send Email"
}

func (a *SendEmailAction) Description() string {
	return "Sends an email to one or more recipients using a customizable template and attachments."
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
	return nil
}

func (a *SendEmailAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"subject": autoform.NewShortTextField().
			SetDisplayName(" Subject").
			SetDescription("The subject of the message").
			SetRequired(true).
			Build(),
		"to": autoform.NewShortTextField().
			SetDisplayName(" To").
			SetDescription("The receiver of the message address").
			SetRequired(true).
			Build(),
		"body": autoform.NewLongTextField().
			SetDisplayName("Body").
			SetDescription("Email body").
			SetRequired(true).
			Build(),
		"bcc": autoform.NewShortTextField().
			SetDisplayName("BCC").
			SetDescription("Blind carbon copy recipients (comma-separated)").
			SetRequired(false).
			Build(),
		"cc": autoform.NewLongTextField().
			SetDisplayName("CC").
			SetDescription("Carbon copy recipients (comma-separated)").
			SetRequired(false).
			Build(),
	}
}

func (a *SendEmailAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[sendEmailActionProps](ctx.BaseContext)
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

	return "Message sent successfully!", nil
}

func (a *SendEmailAction) Auth() *sdk.Auth {
	return nil
}

func (a *SendEmailAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *SendEmailAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewSendEmailAction() sdk.Action {
	return &SendEmailAction{}
}
