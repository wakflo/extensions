package actions

import (
	"context"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type getThreadActionProps struct {
	ThreadID string `json:"id"`
	Format   string `json:"format"`
}

type GetThreadAction struct{}

// Metadata returns metadata about the action
func (a *GetThreadAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_thread",
		DisplayName:   "Get Thread",
		Description:   "Retrieves a specific thread or conversation from a messaging platform, allowing you to incorporate its contents into your automated workflow.",
		Type:          core.ActionTypeAction,
		Documentation: getThreadDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GetThreadAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_thread", "Get Thread")

	form.TextField("id", "id").
		Placeholder("Thread ID").
		HelpText("The thread Id of the mail to read").
		Required(true)

	form.SelectField("format", "format").
		Placeholder("Format").
		HelpText("The format of the email to read").
		AddOption("full", "Full").
		AddOption("metadata", "Metadata").
		AddOption("minimal", "Minimal").
		AddOption("raw", "Raw").
		Required(false)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetThreadAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GetThreadAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getThreadActionProps](ctx)
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

func NewGetThreadAction() sdk.Action {
	return &GetThreadAction{}
}
