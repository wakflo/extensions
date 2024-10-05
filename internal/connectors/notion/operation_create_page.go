package notion

import (
	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createPageProps struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	DatabaseID string `json:"database"`
	PageID     string `json:"page_id"`
}

type CreatePageOperation struct {
	options *sdk.OperationInfo
}

func NewCreatePageOperation() *CreatePageOperation {
	return &CreatePageOperation{
		options: &sdk.OperationInfo{
			Name:        "Create a Notion Page",
			Description: "Create a new page in a Notion database",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"title": autoform.NewShortTextField().
					SetDisplayName("Title").
					SetDescription("The title of the page").
					SetRequired(true).
					Build(),
				"content": autoform.NewLongTextField().
					SetDisplayName("Content").
					SetDescription("The content of the page").
					SetRequired(true).
					Build(),
				"database": getNotionDatabasesInput("Database", "Select a database", true),
				"page_id":  getNotionPagesInput("Page", "Select a page", false),
			},
		},
	}
}
