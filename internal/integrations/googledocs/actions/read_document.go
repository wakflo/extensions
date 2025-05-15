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

type readDocumentActionProps struct {
	DocumentID string `json:"id"`
}

type ReadDocumentAction struct{}

// Metadata returns metadata about the action
func (a *ReadDocumentAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "read_document",
		DisplayName:   "Read Document",
		Description:   "Reads and extracts data from a document, such as a PDF or Word file, into a workflow. This action can be used to retrieve specific information, extract text, or parse structured data from a document. The extracted data can then be used in subsequent actions within the workflow.",
		Type:          core.ActionTypeAction,
		Documentation: readDocumentDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ReadDocumentAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("read_document", "Read Document")

	form.TextField("id", "Document ID").
		Placeholder("Document ID").
		HelpText("The id of the document.").
		Required(true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ReadDocumentAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ReadDocumentAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[readDocumentActionProps](ctx)
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

func NewReadDocumentAction() sdk.Action {
	return &ReadDocumentAction{}
}
