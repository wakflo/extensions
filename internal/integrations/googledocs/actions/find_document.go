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

type findDocumentActionProps struct {
	DocumentID string `json:"id"`
}

type FindDocumentAction struct{}

func (a *FindDocumentAction) Name() string {
	return "Find Document"
}

func (a *FindDocumentAction) Description() string {
	return "Searches for and retrieves a specific document from a designated repository or database, allowing you to automate tasks that require access to existing documents."
}

func (a *FindDocumentAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *FindDocumentAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &findDocumentDocs,
	}
}

func (a *FindDocumentAction) Icon() *string {
	return nil
}

func (a *FindDocumentAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"id": autoform.NewShortTextField().
			SetDisplayName("Document ID").
			SetDescription("The id of the document.").
			SetRequired(true).
			Build(),
	}
}

func (a *FindDocumentAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[findDocumentActionProps](ctx.BaseContext)
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

func (a *FindDocumentAction) Auth() *sdk.Auth {
	return nil
}

func (a *FindDocumentAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *FindDocumentAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewFindDocumentAction() sdk.Action {
	return &FindDocumentAction{}
}
