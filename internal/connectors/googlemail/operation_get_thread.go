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

type getThreadOperationProps struct {
	ThreadID string `json:"id"`
	Format   string `json:"format"`
}

type GetThreadOperation struct {
	options *sdk.OperationInfo
}

func NewGetThreadOperation() *GetThreadOperation {
	return &GetThreadOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Thread",
			Description: "Get a thread from your Gmail account via Id",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"id": autoform.NewShortTextField().
					SetDisplayName("Thread ID").
					SetDescription("The thread Id of the mail to read").
					SetRequired(true).
					Build(),
				"format": autoform.NewSelectField().
					SetDisplayName("Format").
					SetDescription("The format of the email to read").
					SetOptions(viewMailFormat).
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

func (c *GetThreadOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[getThreadOperationProps](ctx)
	gmailService, err := gmail.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	chosenFormat := input.Format
	if chosenFormat == "" {
		chosenFormat = "full"
	}

	mail, err := gmailService.Users.Messages.Get("me", input.ThreadID).
		Format(chosenFormat).
		Do()
	if err != nil {
		return nil, err
	}

	return mail, nil
}

func (c *GetThreadOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetThreadOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
