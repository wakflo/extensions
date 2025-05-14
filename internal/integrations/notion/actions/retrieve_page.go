package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/notion/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type retrievePageActionProps struct {
	DatabaseID string `json:"databaseId"`
	PageID     string `json:"page_id"`
}

type RetrievePageAction struct{}

// func (a *RetrievePageAction) Name() string {
// 	return "Retrieve Page"
// }

func (a *RetrievePageAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "retrieve_page",
		DisplayName:   "Retrieve Page",
		Description:   "Retrieve Page: Extracts and retrieves the content of a specified web page, allowing you to access and analyze the page's text, images, or other relevant information.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: retrievePageDocs,
		SampleOutput: map[string]any{
			"object": "page",
			"id":     "1234567890",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *RetrievePageAction) Properties() *smartform.FormSchema {

	form := smartform.NewForm("retrieve_page", "Retrieve Page")

	shared.GetNotionDatabasesProp(form)

	shared.GetNotionPagesProp("Page ID", "Select a page", false, form)

	schema := form.Build()

	return schema
}

func (a *RetrievePageAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[retrievePageActionProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	if input.PageID == "" {
		return nil, errors.New("page ID is required")
	}

	pageData, err := shared.GetNotionPage(token, input.PageID)
	if err != nil {
		return nil, err
	}

	return pageData, nil
}

func (a *RetrievePageAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewRetrievePageAction() sdk.Action {
	return &RetrievePageAction{}
}
