package notion

import (
	"errors"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type retrievePageOperationProps struct {
	DatabaseID string `json:"database"`
	PageID     string `json:"page_id"`
}

type RetrievePageOperation struct {
	options *sdk.OperationInfo
}

func NewRetrievePageOperation() *RetrievePageOperation {
	return &RetrievePageOperation{
		options: &sdk.OperationInfo{
			Name:        "Retrieve a Notion Page",
			Description: "Retrieve a page from Notion using the page ID",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"database": getNotionDatabasesInput("Database", "Select a database", true),
				"page_id":  getNotionPagesInput("Page ID", "Select a page", false),
			},
		},
	}
}

func (r *RetrievePageOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing notion auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[retrievePageOperationProps](ctx)

	if input.PageID == "" {
		return nil, errors.New("page ID is required")
	}

	pageData, err := getNotionPage(accessToken, input.PageID)
	if err != nil {
		return nil, err
	}

	return pageData, nil
}

func (r *RetrievePageOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return r.Run(ctx)
}

func (r *RetrievePageOperation) GetInfo() *sdk.OperationInfo {
	return r.options
}
