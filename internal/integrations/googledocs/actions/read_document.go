package actions

import (
	"context"
	"errors"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

type readDocumentActionProps struct {
	DocumentID string `json:"id"`
}

type ReadDocumentAction struct{}

func (a *ReadDocumentAction) Name() string {
	return "Read Document"
}

func (a *ReadDocumentAction) Description() string {
	return "Reads and extracts data from a document, such as a PDF or Word file, into a workflow. This action can be used to retrieve specific information, extract text, or parse structured data from a document. The extracted data can then be used in subsequent actions within the workflow."
}

func (a *ReadDocumentAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ReadDocumentAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &readDocumentDocs,
	}
}

func (a *ReadDocumentAction) Icon() *string {
	return nil
}

func (a *ReadDocumentAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"id": autoform.NewShortTextField().
			SetDisplayName("Document ID").
			SetDescription("The id of the document.").
			SetRequired(true).
			Build(),
	}
}

func (a *ReadDocumentAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[readDocumentActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	docService, err := docs.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	if input.DocumentID == "" {
		return nil, errors.New("name is required")
	}

	document, err := docService.Documents.Get(input.DocumentID).
		Do()
	return document, err
}

func (a *ReadDocumentAction) Auth() *sdk.Auth {
	return nil
}

func (a *ReadDocumentAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *ReadDocumentAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewReadDocumentAction() sdk.Action {
	return &ReadDocumentAction{}
}
