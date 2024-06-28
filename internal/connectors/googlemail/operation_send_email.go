package googlemail

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type sendMailOperationProps struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	CC      string `json:"cc"`
	BCC     string `json:"bcc"`
}

type SendMailOperation struct {
	options *sdk.OperationInfo
}

func NewSendMailOperation() sdk.IOperation {
	return &SendMailOperation{
		options: &sdk.OperationInfo{
			Name:        "Send Email",
			Description: "Send an email from your Gmail account",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
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
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *SendMailOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[sendMailOperationProps](ctx)
	gmailService, err := gmail.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	userProfile, err := gmailService.Users.GetProfile("me").Do()
	if err != nil {
		return nil, err
	}

	fromEmail := userProfile.EmailAddress

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

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func (c *SendMailOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *SendMailOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
