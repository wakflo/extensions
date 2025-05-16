package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/notion/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type createPageActionProps struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	DatabaseID string `json:"database"`
	PageID     string `json:"page_id"`
}

type CreatePageAction struct{}

func (a *CreatePageAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_page",
		DisplayName:   "Create Page",
		Description:   "Create Page: Automatically generates a new page within your website or application, allowing you to quickly create and deploy new content without manual intervention.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: createPageDocs,
		SampleOutput: map[string]any{
			"object": "page",
			"id":     "1234567890",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *CreatePageAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_page", "Create Page")

	form.TextField("title", "Title").
		Placeholder("Enter a title").
		Required(true).
		HelpText("The title of the page")

	form.TextareaField("content", "Content").
		Placeholder("Enter content").
		Required(true).
		HelpText("The content of the page")

	shared.GetNotionDatabasesProp(form)

	schema := form.Build()

	return schema
}

func (a *CreatePageAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createPageActionProps](ctx)
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

	notionPage, _ := shared.CreateNotionPage(token, input.PageID, input.Title, input.Content)

	return notionPage, nil
}

func (a *CreatePageAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewCreatePageAction() sdk.Action {
	return &CreatePageAction{}
}
