package actions

import (
	"errors"

	"github.com/wakflo/extensions/internal/integrations/notion/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type retrievePageActionProps struct {
	DatabaseID string `json:"databaseId"`
	PageID     string `json:"page_id"`
}

type RetrievePageAction struct{}

func (a *RetrievePageAction) Name() string {
	return "Retrieve Page"
}

func (a *RetrievePageAction) Description() string {
	return "Retrieves the content of a specified web page and extracts relevant information, such as text, images, or links, to be used in subsequent workflow steps."
}

func (a *RetrievePageAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *RetrievePageAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &retrievePageDocs,
	}
}

func (a *RetrievePageAction) Icon() *string {
	return nil
}

func (a *RetrievePageAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"databaseId": shared.GetNotionDatabasesInput("Database", "Select a database", true),
		"page_id":    shared.GetNotionPagesInput("Page ID", "Select a page", false),
	}
}

func (a *RetrievePageAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[retrievePageActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing calendly auth token")
	}
	accessToken := ctx.Auth.AccessToken

	if input.PageID == "" {
		return nil, errors.New("page ID is required")
	}

	pageData, err := shared.GetNotionPage(accessToken, input.PageID)
	if err != nil {
		return nil, err
	}

	return pageData, nil
}

func (a *RetrievePageAction) Auth() *sdk.Auth {
	return nil
}

func (a *RetrievePageAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *RetrievePageAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewRetrievePageAction() sdk.Action {
	return &RetrievePageAction{}
}
