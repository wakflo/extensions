package actions

import (
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/facebookpages/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type deletePostActionProps struct {
	PageID string `json:"page_id"`
	PostID string `json:"post_id"`
}

type DeletePostAction struct{}

// Metadata returns metadata about the action
func (c *DeletePostAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "delete_post",
		DisplayName:   "Delete a Post",
		Description:   "Delete a post from a page.",
		Type:          core.ActionTypeAction,
		Documentation: deletePostDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"success": true,
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (c *DeletePostAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("delete_post", "Delete a Post")

	shared.RegisterFacebookPageProps(form)

	shared.RegisterPagePostsProps(form)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (c *DeletePostAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (c *DeletePostAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[deletePostActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	url := "/" + input.PostID

	pageAccessToken, err := shared.GetPageAccessToken(authCtx.Token.AccessToken, input.PageID)
	if err != nil {
		return nil, err
	}
	deletedPost, err := shared.PostActionFunc(pageAccessToken, http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return deletedPost, nil
}

func NewDeletePostAction() sdk.Action {
	return &DeletePostAction{}
}
