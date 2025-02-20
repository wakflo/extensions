package actions

import (
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/facebookpages/shared"
	sdkcore "github.com/wakflo/go-sdk/core"

	"github.com/wakflo/go-sdk/sdk"
)

type getPostActionProps struct {
	PageID string `json:"page_id"`
}

type GetPostsAction struct{}

func (c *GetPostsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (c GetPostsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c GetPostsAction) Name() string {
	return "Get facebook Posts"
}

func (c GetPostsAction) Description() string {
	return "get posts."
}

func (c GetPostsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getPostsDocs,
	}
}

func (c GetPostsAction) Icon() *string {
	return nil
}

func (c GetPostsAction) SampleData() sdkcore.JSON {
	return nil
}

func (c GetPostsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"page_id": shared.GetFacebookPageInput("Select a Page", "The ID of the page you want to get the post from", true),
	}
}

func (c GetPostsAction) Auth() *sdk.Auth {
	return &sdk.Auth{
		Inherit: true,
	}
}

func (c GetPostsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getPostActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	pageAccessToken, err := shared.GetPageAccessToken(ctx.Auth.AccessToken, input.PageID)
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
