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

type createDocumentOperationProps struct {
	Name string `json:"name"`
	// Body any `json:"body"`
}

type CreateDocumentOperation struct {
	options *sdk.OperationInfo
}

func NewCreateDocumentOperation() *CreateDocumentOperation {
	return &CreateDocumentOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Document",
			Description: "Search for document by name",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"name": autoform.NewShortTextField().
					SetDisplayName("Document Name").
					SetDescription("The name of the document.").
					SetRequired(true).
					Build(),
				// "body": autoform.NewLongTextField().
				// SetDisplayName("Document Body").
				// SetDescription("The initial content of the document.").
				// SetRequired(true).
				// Build(),

				// "name": autoform.NewShortTextField().
				// SetDisplayName("Document Name").
				// SetDescription("The name of the new folder.").
				// SetRequired(true).
				// Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *CreateDocumentOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[createDocumentOperationProps](ctx)
	docService, err := docs.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	if input.Name == "" {
		return nil, errors.New("name is required")
	}

	// if input.Body == "" {
	// 	return nil, errors.New("name is required")
	// }

	document, err := docService.Documents.Create(&docs.Document{
		Title: input.Name,
		// Body: input.Body,
	}).
		// Fields("id, name, mimeType, webViewLink, kind, createdTime").
		Do()
	return document, err
}

func (c *CreateDocumentOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateDocumentOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
