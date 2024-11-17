package notion

import (
	"errors"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type updatePageOperationProps struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	DatabaseID string `json:"database"`
	PageID     string `json:"page_id"`
}

type UpdatePageOperation struct {
	options *sdk.OperationInfo
}

func NewUpdatePageOperation() *UpdatePageOperation {
	return &UpdatePageOperation{
		options: &sdk.OperationInfo{
			Name:        "Update a Notion Page",
			Description: "Update a page in a Notion database",
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
				"page_id":  getNotionPagesInput("Page", "Select a page", true),
			},
		},
	}
}

func (c *UpdatePageOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing calendly auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[updatePageOperationProps](ctx)

	if input.Title == "" {
		return nil, errors.New("title is required")
	}

	if input.DatabaseID == "" {
		return nil, errors.New("database is required")
	}

	if input.PageID == "" {
		return nil, errors.New("parent page is required")
	}

	notionPage, _ := updateNotionPage(accessToken, input.PageID, input.Title, input.Content)

	return notionPage, nil
}

func (c *UpdatePageOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *UpdatePageOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
