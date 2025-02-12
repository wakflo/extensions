package actions

import (
	"context"
	"errors"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type listMailsActionProps struct {
	Label string `json:"label"`
}

type ListMailsAction struct{}

func (a *ListMailsAction) Name() string {
	return "List Mails"
}

func (a *ListMailsAction) Description() string {
	return "Retrieve a list of emails from your email account or service, allowing you to automate workflows based on specific mail criteria."
}

func (a *ListMailsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListMailsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listMailsDocs,
	}
}

func (a *ListMailsAction) Icon() *string {
	return nil
}

func (a *ListMailsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"label": autoform.NewShortTextField().
			SetDisplayName("Label").
			SetDescription("The mail label to read from (e.g, inbox, sent, drafts, etc)").
			SetRequired(true).
			Build(),
	}
}

func (a *ListMailsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listMailsActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

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

func (a *ListMailsAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListMailsAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *ListMailsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListMailsAction() sdk.Action {
	return &ListMailsAction{}
}
