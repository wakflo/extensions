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

type listMailsActionProps struct {
	Label string `json:"label"`
}

type ListMailsAction struct{}

// Metadata returns metadata about the action
func (a *ListMailsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_mails",
		DisplayName:   "List Mails",
		Description:   "Retrieve a list of emails from your email account or service, allowing you to automate workflows based on specific mail criteria.",
		Type:          core.ActionTypeAction,
		Documentation: listMailsDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ListMailsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_mails", "List Mails")

	form.TextField("label", "label").
		Placeholder("Label").
		HelpText("The mail label to read from (e.g, inbox, sent, drafts, etc)").
		Required(true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ListMailsAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ListMailsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[listMailsActionProps](ctx)
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

func NewListMailsAction() sdk.Action {
	return &ListMailsAction{}
}
