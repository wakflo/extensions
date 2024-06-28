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

type findDocumentOperationProps struct {
	DocumentID string `json:"id"`
}

type FindDocumentOperation struct {
	options *sdk.OperationInfo
}

func NewFindDocumentOperation() *FindDocumentOperation {
	return &FindDocumentOperation{
		options: &sdk.OperationInfo{
			Name:        "Find Document",
			Description: "Search for document by ID.",
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

func (c *FindDocumentOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[findDocumentOperationProps](ctx)
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

func (c *FindDocumentOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *FindDocumentOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
