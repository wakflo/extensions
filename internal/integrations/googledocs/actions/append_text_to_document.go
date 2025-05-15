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

type appendTextToDocumentActionProps struct {
	DocumentID string `json:"id"`
	Text       string `json:"text"`
}

type AppendTextToDocumentAction struct{}

// Metadata returns metadata about the action
func (a *AppendTextToDocumentAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "append_text_to_document",
		DisplayName:   "Append Text to Document",
		Description:   "Appends text to an existing document or file, allowing you to add notes, comments, or additional information to the end of the file. This integration action is useful when you need to update a document with new information, such as adding a signature or timestamp, without modifying the original content.",
		Type:          core.ActionTypeAction,
		Documentation: appendTextToDocumentDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *AppendTextToDocumentAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("append_text_to_document", "Append Text to Document")

	form.TextField("id", "Document ID").
		Placeholder("Document ID").
		HelpText("The id of the document.").
		Required(true)

	form.TextareaField("text", "Text").
		Placeholder("Text to append").
		HelpText("The text to append to the document").
		Required(true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *AppendTextToDocumentAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *AppendTextToDocumentAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[appendTextToDocumentActionProps](ctx)
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
		return nil, errors.New("id is required")
	}

	if input.Text == "" {
		return nil, errors.New("text is required")
	}

	document, err := docService.Documents.BatchUpdate(input.DocumentID, &docs.BatchUpdateDocumentRequest{
		Requests: []*docs.Request{
			{
				InsertText: &docs.InsertTextRequest{
					Text:                 input.Text,
					EndOfSegmentLocation: &docs.EndOfSegmentLocation{},
				},
			},
		},
	}).Do()

	return document, err
}

func NewAppendTextToDocumentAction() sdk.Action {
	return &AppendTextToDocumentAction{}
}
