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

type findDocumentActionProps struct {
	DocumentID string `json:"id"`
}

type FindDocumentAction struct{}

// Metadata returns metadata about the action
func (a *FindDocumentAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "find_document",
		DisplayName:   "Find Document",
		Description:   "Searches for and retrieves a specific document from a designated repository or database, allowing you to automate tasks that require access to existing documents.",
		Type:          core.ActionTypeAction,
		Documentation: findDocumentDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *FindDocumentAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("find_document", "Find Document")

	form.TextField("id", "Document ID").
		Placeholder("Document ID").
		HelpText("The id of the document.").
		Required(true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *FindDocumentAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *FindDocumentAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[findDocumentActionProps](ctx)
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

	if input.DocumentID == "" {
		return nil, errors.New("document ID is required")
	}

	document, err := docService.Documents.Get(input.DocumentID).
		Do()
	return document, err
}

func NewFindDocumentAction() sdk.Action {
	return &FindDocumentAction{}
}
