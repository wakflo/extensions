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

type createDocumentActionProps struct {
	Name string `json:"name"`
	// Body any `json:"body"`
}

type CreateDocumentAction struct{}

func (a *CreateDocumentAction) Name() string {
	return "Create Document"
}

func (a *CreateDocumentAction) Description() string {
	return "Create Document: Automatically generates and saves a new document based on pre-defined templates and data inputs, allowing you to streamline your document creation process and reduce manual errors."
}

func (a *CreateDocumentAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateDocumentAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createDocumentDocs,
	}
}

func (a *CreateDocumentAction) Icon() *string {
	return nil
}

func (a *CreateDocumentAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"name": autoform.NewShortTextField().
			SetDisplayName("Document Name").
			SetDescription("The name of the document.").
			SetRequired(true).
			Build(),
	}
}

func (a *CreateDocumentAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createDocumentActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	docService, err := docs.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
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

func (a *CreateDocumentAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateDocumentAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateDocumentAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateDocumentAction() sdk.Action {
	return &CreateDocumentAction{}
}
