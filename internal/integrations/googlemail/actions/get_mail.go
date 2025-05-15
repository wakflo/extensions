package actions

import (
	"context"
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type getMailActionProps struct {
	MailID string `json:"id"`
}

type GetMailAction struct{}

// Metadata returns metadata about the action
func (a *GetMailAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_mail",
		DisplayName:   "Get Mail",
		Description:   "Retrieves emails from a specified email account or inbox, allowing you to automate tasks triggered by new mail arrivals.",
		Type:          core.ActionTypeAction,
		Documentation: getMailDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GetMailAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_mail", "Get Mail")

	form.TextField("id", "id").
		Placeholder("Message ID").
		HelpText("The messageId of the mail to read").
		Required(true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetMailAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GetMailAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getMailActionProps](ctx)
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

func NewGetMailAction() sdk.Action {
	return &GetMailAction{}
}
