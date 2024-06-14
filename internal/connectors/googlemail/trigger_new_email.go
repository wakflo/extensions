package googlemail

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type triggerNewEmailProps struct {
	Subject        string     `json:"subject"`
	From           string     `json:"from"`
	MaxResults     int        `json:"maxResults"`
	RecievedTime   *time.Time `json:"receivedTime"`
	RecievedTimeOp *string    `json:"receivedTimeOp"`
}

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
			ErrorSettings: sdkcore.StepErrorSettings{
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

	input := sdk.InputToType[triggerNewEmailProps](ctx)
	gmailService, err := gmail.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	var qarr []string
	if input.Subject != "" {
		qarr = append(qarr, fmt.Sprintf("subject:%v", input.Subject))
	}
	if input.From != "" {
		qarr = append(qarr, fmt.Sprintf("from:%v", input.From))
	}
	if input.RecievedTime != nil {
		input.RecievedTime = ctx.LastRun
	}
	if input.RecievedTime != nil {
		op := ">"
		if input.RecievedTimeOp != nil {
			op = *input.RecievedTimeOp
		}
		qarr = append(qarr, fmt.Sprintf(`received time %v '%v'`, op, input.RecievedTime.UTC().Format("2006-01-02T15:04:05Z")))
	}

	q := strings.Join(qarr, " ")

	req := gmailService.Users.Messages.List("me").
		Q(q).MaxResults(int64(input.MaxResults))

	res, err := req.Do()
	if err != nil {
		return nil, err
	}

	return res.Messages, nil
}

func (t TriggerNewEmail) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return t.Run(ctx)
}

func (t TriggerNewEmail) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t TriggerNewEmail) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t TriggerNewEmail) GetInfo() *sdk.TriggerInfo {
	return t.options
}
