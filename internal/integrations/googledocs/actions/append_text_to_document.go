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

type appendTextToDocumentActionProps struct {
	DocumentID string `json:"id"`
	Text       string `json:"text"`
}

type AppendTextToDocumentAction struct{}

func (a *AppendTextToDocumentAction) Name() string {
	return "Append Text to Document"
}

func (a *AppendTextToDocumentAction) Description() string {
	return "Appends text to an existing document or file, allowing you to add notes, comments, or additional information to the end of the file. This integration action is useful when you need to update a document with new information, such as adding a signature or timestamp, without modifying the original content."
}

func (a *AppendTextToDocumentAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *AppendTextToDocumentAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &appendTextToDocumentDocs,
	}
}

func (a *AppendTextToDocumentAction) Icon() *string {
	return nil
}

func (a *AppendTextToDocumentAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"id": autoform.NewShortTextField().
			SetDisplayName("Document ID").
			SetDescription("The id of the document.").
			SetRequired(true).
			Build(),
		"text": autoform.NewLongTextField().
			SetDisplayName("Text to append").
			SetDescription("The text to append to the document").
			SetRequired(true).
			Build(),
	}
}

func (a *AppendTextToDocumentAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[appendTextToDocumentActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	docService, err := docs.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
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
	}).
		// Fields("id, name, mimeType, webViewLink, kind, createdTime").
		Do()
	return document, err
}

func (a *AppendTextToDocumentAction) Auth() *sdk.Auth {
	return nil
}

func (a *AppendTextToDocumentAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *AppendTextToDocumentAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewAppendTextToDocumentAction() sdk.Action {
	return &AppendTextToDocumentAction{}
}
