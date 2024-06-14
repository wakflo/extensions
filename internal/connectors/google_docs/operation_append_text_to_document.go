package googledocs

import (
	"context"
	"errors"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type appendTextToDocumentOperationProps struct {
	DocumentID string `JSON:"id"`
	Text       string `JSON:"text"`
}

type AppendTextToDocumentOperation struct {
	options *sdk.OperationInfo
}

func NewAppendTextToDocumentOperation() *AppendTextToDocumentOperation {
	return &AppendTextToDocumentOperation{
		options: &sdk.OperationInfo{
			Name:        "Append Text",
			Description: "Appends text to google docs",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
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
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *AppendTextToDocumentOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[appendTextToDocumentOperationProps](ctx)
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

func (c *AppendTextToDocumentOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *AppendTextToDocumentOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
