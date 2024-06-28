package googlemail

import (
	"context"
	"errors"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type getMailByIDOperationProps struct {
	MailID string `json:"id"`
}

type GetMailByIDOperation struct {
	options *sdk.OperationInfo
}

func NewGetMailByIDOperation() *GetMailByIDOperation {
	return &GetMailByIDOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Email",
			Description: "Get an email from your Gmail account via Id",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"id": autoform.NewShortTextField().
					SetDisplayName("Message ID").
					SetDescription("The messageId of the mail to read").
					SetRequired(true).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *GetMailByIDOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[getMailByIDOperationProps](ctx)
	gmailService, err := gmail.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	if input.MailID == "" {
		return nil, errors.New("mail ID is required")
	}

	mail, err := gmailService.Users.Messages.Get("me", input.MailID).
		Format("full").
		Do()
	if err != nil {
		return nil, err
	}

	return mail, nil
}

func (c *GetMailByIDOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetMailByIDOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
