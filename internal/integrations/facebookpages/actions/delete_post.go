package actions

import (
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/facebookpages/shared"
	sdkcore "github.com/wakflo/go-sdk/core"

	"github.com/wakflo/go-sdk/sdk"
)

type deletePostActionProps struct {
	PageID string `json:"page_id"`
	PostID string `json:"post_id"`
}

type DeletePostAction struct{}

func (c *DeletePostAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (c DeletePostAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c DeletePostAction) Name() string {
	return "Delete a Post"
}

func (c DeletePostAction) Description() string {
	return "delete a post from a page."
}

func (c DeletePostAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &deletePostDocs,
	}
}

func (c DeletePostAction) Icon() *string {
	return nil
}

func (c DeletePostAction) SampleData() sdkcore.JSON {
	return nil
}

func (c DeletePostAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"page_id": shared.GetFacebookPageInput("Select a page", "The page you want to get the post from", true),
		"post_id": shared.GetPagePostsInput("Select a post", "The post you want to update", true),
	}
}

func (c DeletePostAction) Auth() *sdk.Auth {
	return &sdk.Auth{
		Inherit: true,
	}
}

func (c DeletePostAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[deletePostActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	url := "/" + input.PostID

	pageAccessToken, err := shared.GetPageAccessToken(ctx.Auth.AccessToken, input.PageID)
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
