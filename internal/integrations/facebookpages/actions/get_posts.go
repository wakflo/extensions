package actions

import (
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/facebookpages/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getPostActionProps struct {
	PageID string `json:"page_id"`
}

type GetPostsAction struct{}

// Metadata returns metadata about the action
func (c *GetPostsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_facebook_posts",
		DisplayName:   "Get Facebook Posts",
		Description:   "Get posts from a Facebook page.",
		Type:          core.ActionTypeAction,
		Documentation: getPostsDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"data": []map[string]any{
				{
					"id":           "12345678901234567_98765432109876543",
					"message":      "Example post content",
					"created_time": "2023-03-15T12:34:56+0000",
				},
			},
			"paging": map[string]any{
				"cursors": map[string]string{
					"before": "example_before_cursor",
					"after":  "example_after_cursor",
				},
				"next": "https://graph.facebook.com/v16.0/next_page",
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (c *GetPostsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_facebook_posts", "Get Facebook Posts")

	shared.RegisterFacebookPageProps(form)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (c *GetPostsAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (c *GetPostsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getPostActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	pageAccessToken, err := shared.GetPageAccessToken(authCtx.AccessToken, input.PageID)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("/%s/feed", input.PageID)

	posts, err := shared.ActionFunc(pageAccessToken, url, nil)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func NewGetPostsAction() sdk.Action {
	return &GetPostsAction{}
}
