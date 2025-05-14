package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/notion/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type updatePageActionProps struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	DatabaseID string `json:"database"`
	PageID     string `json:"page_id"`
}

type UpdatePageAction struct{}

func (a *UpdatePageAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "update_page",
		DisplayName:   "Update Page",
		Description:   "Update Page: Updates an existing page with new content, replacing any existing text, images, or other elements.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: updatePageDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *UpdatePageAction) Properties() *smartform.FormSchema {

	form := smartform.NewForm("update_page", "Update Page")

	form.TextField("title", "Title").
		Placeholder("Enter a title").
		Required(true).
		HelpText("The title of the page")

	form.TextareaField("content", "Content").
		Placeholder("Enter content").
		Required(true).
		HelpText("The content of the page")

	shared.GetNotionDatabasesProp(form)

	shared.GetNotionPagesProp("Page ID", "Select a page", false, form)

	schema := form.Build()

	return schema
}

func (a *UpdatePageAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updatePageActionProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	if input.Title == "" {
		return nil, errors.New("title is required")
	}

	if input.DatabaseID == "" {
		return nil, errors.New("database is required")
	}

	if input.PageID == "" {
		return nil, errors.New("parent page is required")
	}

	notionPage, _ := shared.UpdateNotionPage(token, input.PageID, input.Title, input.Content)

	return notionPage, nil
}

func (a *UpdatePageAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewUpdatePageAction() sdk.Action {
	return &UpdatePageAction{}
}
