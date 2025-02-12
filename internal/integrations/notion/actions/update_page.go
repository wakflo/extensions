package actions

import (
	"errors"

	"github.com/wakflo/extensions/internal/integrations/notion/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type updatePageActionProps struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	DatabaseID string `json:"database"`
	PageID     string `json:"page_id"`
}

type UpdatePageAction struct{}

func (a *UpdatePageAction) Name() string {
	return "Update Page"
}

func (a *UpdatePageAction) Description() string {
	return "Updates an existing page with new content, replacing any existing text, images, or other elements."
}

func (a *UpdatePageAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *UpdatePageAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &updatePageDocs,
	}
}

func (a *UpdatePageAction) Icon() *string {
	return nil
}

func (a *UpdatePageAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
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
		"database": shared.GetNotionDatabasesInput("Database", "Select a database", true),
		"page_id":  shared.GetNotionPagesInput("Page", "Select a page", true),
	}
}

func (a *UpdatePageAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updatePageActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing calendly auth token")
	}
	accessToken := ctx.Auth.AccessToken

	if input.Title == "" {
		return nil, errors.New("title is required")
	}

	if input.DatabaseID == "" {
		return nil, errors.New("database is required")
	}

	if input.PageID == "" {
		return nil, errors.New("parent page is required")
	}

	notionPage, _ := shared.UpdateNotionPage(accessToken, input.PageID, input.Title, input.Content)

	return notionPage, nil
}

func (a *UpdatePageAction) Auth() *sdk.Auth {
	return nil
}

func (a *UpdatePageAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *UpdatePageAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewUpdatePageAction() sdk.Action {
	return &UpdatePageAction{}
}
