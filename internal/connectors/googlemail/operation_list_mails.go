package googlemail

import (
	"context"
	"errors"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type listMailsOperationProps struct {
	Label string `json:"label"`
}

type ListMailsOperation struct {
	options *sdk.OperationInfo
}

func NewListMailsOperation() sdk.IOperation {
	return &ListMailsOperation{
		options: &sdk.OperationInfo{
			Name:        "List Email",
			Description: "Get a list of emails from your Gmail account",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"label": autoform.NewShortTextField().
					SetDisplayName("Label").
					SetDescription("The mail label to read from (e.g, inbox, sent, drafts, etc)").
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

func (c *ListMailsOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[listMailsOperationProps](ctx)
	gmailService, err := gmail.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	if input.Label == "" {
		return nil, errors.New("label name is required")
	}

	query := "in:" + input.Label

	listResponse, err := gmailService.Users.Messages.List("me").Q(query).Do()
	if err != nil {
		return nil, err
	}

	return listResponse, nil
}

func (c *ListMailsOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *ListMailsOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
