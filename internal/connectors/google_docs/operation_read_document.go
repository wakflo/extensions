package googledocs

import (
	"context"
	"errors"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

type readDocumentOperationProps struct {
	DocumentID string `json:"id"`
}

type ReadDocumentOperation struct {
	options *sdk.OperationInfo
}

func NewReadDocumentOperation() *ReadDocumentOperation {
	return &ReadDocumentOperation{
		options: &sdk.OperationInfo{
			Name:        "Read Document",
			Description: "Read document by ID.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"id": autoform.NewShortTextField().
					SetDisplayName("Document ID").
					SetDescription("The id of the document.").
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

func (c *ReadDocumentOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[readDocumentOperationProps](ctx)
	docService, err := docs.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	if input.DocumentID == "" {
		return nil, errors.New("name is required")
	}

	document, err := docService.Documents.Get(input.DocumentID).
		// Fields("id, name, mimeType, webViewLink, kind, createdTime").
		Do()
	return document, err
}

func (c *ReadDocumentOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *ReadDocumentOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
