package actions

import (
	"context"
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

type createDocumentActionProps struct {
	Name string `json:"name"`
	// Body any `json:"body"`
}

type CreateDocumentAction struct{}

// Metadata returns metadata about the action
func (a *CreateDocumentAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_document",
		DisplayName:   "Create Document",
		Description:   "Create Document: Automatically generates and saves a new document based on pre-defined templates and data inputs, allowing you to streamline your document creation process and reduce manual errors.",
		Type:          core.ActionTypeAction,
		Documentation: createDocumentDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CreateDocumentAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_document", "Create Document")

	form.TextField("name", "Name").
		Placeholder("Document Name").
		HelpText("The name of the document.").
		Required(true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CreateDocumentAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CreateDocumentAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createDocumentActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	docService, err := docs.NewService(context.Background(), option.WithTokenSource(*authCtx.TokenSource))
	if err != nil {
		return nil, err
	}

	if input.Name == "" {
		return nil, errors.New("name is required")
	}

	document, err := docService.Documents.Create(&docs.Document{
		Title: input.Name,
	}).
		Do()
	return document, err
}

func NewCreateDocumentAction() sdk.Action {
	return &CreateDocumentAction{}
}
