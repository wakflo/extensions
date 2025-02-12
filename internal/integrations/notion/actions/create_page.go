package actions

import (
	"errors"

	"github.com/wakflo/extensions/internal/integrations/notion/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createPageActionProps struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	DatabaseID string `json:"database"`
	PageID     string `json:"page_id"`
}

type CreatePageAction struct{}

func (a *CreatePageAction) Name() string {
	return "Create Page"
}

func (a *CreatePageAction) Description() string {
	return "Create Page: Automatically generates a new page within your website or application, allowing you to quickly create and deploy new content without manual intervention."
}

func (a *CreatePageAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreatePageAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createPageDocs,
	}
}

func (a *CreatePageAction) Icon() *string {
	return nil
}

func (a *CreatePageAction) Properties() map[string]*sdkcore.AutoFormSchema {
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
		"page_id":  shared.GetNotionPagesInput("Parent Page", "Select a parent page", false),
	}
}

func (a *CreatePageAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createPageActionProps](ctx.BaseContext)
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

	notionPage, _ := shared.CreateNotionPage(accessToken, input.PageID, input.Title, input.Content)

	return notionPage, nil
}

func (a *CreatePageAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreatePageAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreatePageAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreatePageAction() sdk.Action {
	return &CreatePageAction{}
}
