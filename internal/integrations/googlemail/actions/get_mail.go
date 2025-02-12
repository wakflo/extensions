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

type getMailActionProps struct {
	MailID string `json:"id"`
}

type GetMailAction struct{}

func (a *GetMailAction) Name() string {
	return "Get Mail"
}

func (a *GetMailAction) Description() string {
	return "Retrieves emails from a specified email account or inbox, allowing you to automate tasks triggered by new mail arrivals."
}

func (a *GetMailAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetMailAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getMailDocs,
	}
}

func (a *GetMailAction) Icon() *string {
	return nil
}

func (a *GetMailAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"id": autoform.NewShortTextField().
			SetDisplayName("Message ID").
			SetDescription("The messageId of the mail to read").
			SetRequired(true).
			Build(),
	}
}

func (a *GetMailAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getMailActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

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

func (a *GetMailAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetMailAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *GetMailAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetMailAction() sdk.Action {
	return &GetMailAction{}
}
