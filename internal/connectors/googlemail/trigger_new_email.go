package googlemail

import (
	"context"
	"errors"
	"time"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type TriggerNewEmail struct {
	options *sdk.TriggerInfo
}

func NewTriggerNewEmail() *TriggerNewEmail {
	return &TriggerNewEmail{
		options: &sdk.TriggerInfo{
			Name:        "New Email ",
			Description: "Triggers workflow when new email is received",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"subject": autoform.NewShortTextField().
					SetDisplayName("Subject").
					SetDescription("Subject of the email to trigger on").
					SetRequired(true).
					Build(),
				"from": autoform.NewShortTextField().
					SetDisplayName("From").
					SetDescription("Sender of the email to trigger on").
					SetRequired(false).
					Build(),
				"maxResults": autoform.NewNumberField().
					SetDisplayName("Max Results").
					SetDescription("Maximum number of emails to return").
					SetDefaultValue(50).
					SetRequired(false).
					Build(),
			},
			Type:     sdkcore.TriggerTypeCron,
			Settings: &sdkcore.TriggerSettings{},
			ErrorSettings: &sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (t *TriggerNewEmail) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	gmailService, err := gmail.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	var lastRunTime time.Time
	if ctx.Metadata.LastRun != nil {
		lastRunTime = *ctx.Metadata.LastRun
	}

	query := "after:" + lastRunTime.Format(time.DateOnly)

	messages, err := gmailService.Users.Messages.List("me").Q(query).Do()
	if err != nil {
		return nil, err
	}

	newEmails := make([]map[string]interface{}, 0, len(messages.Messages))
	for _, msg := range messages.Messages {
		email, err := gmailService.Users.Messages.Get("me", msg.Id).Do()
		if err != nil {
			continue
		}

		emailData := map[string]interface{}{
			"id":      email.Id,
			"subject": getHeader(email.Payload.Headers, "Subject"),
			"from":    getHeader(email.Payload.Headers, "From"),
			"date":    getHeader(email.Payload.Headers, "Date"),
		}
		newEmails = append(newEmails, emailData)
	}

	return newEmails, nil
}

func (t *TriggerNewEmail) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return t.Run(ctx)
}

func (t *TriggerNewEmail) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewEmail) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewEmail) GetInfo() *sdk.TriggerInfo {
	return t.options
}
