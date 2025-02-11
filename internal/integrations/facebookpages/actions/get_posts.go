package actions

import (
	"fmt"
	"github.com/wakflo/extensions/internal/integrations/facebookpages/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"

	"github.com/wakflo/go-sdk/integration"
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

func (c GetPostsAction) Documentation() *integration.OperationDocumentation {
	return &integration.OperationDocumentation{
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
		//"page_id": shared.GetFacebookPageInput("Select a Page", "The ID of the page you want to get the post from", true),
		"page_id": autoform.NewShortTextField().
			SetDisplayName("Page").
			SetDescription("Page").
			SetRequired(true).
			Build(),
	}
}

func (c GetPostsAction) Auth() *integration.Auth {
	return &integration.Auth{
		Inherit: true,
	}
}

func (c GetPostsAction) Perform(ctx integration.PerformContext) (sdkcore.JSON, error) {
	input, err := integration.InputToTypeSafely[getPostActionProps](ctx.BaseContext)

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

	//videoResult, err := shared.MakeFacebookRequest(http.MethodGet, ctx.Auth.AccessToken, endpoint, nil)
	//if err != nil {
	//	return nil, err
	//}

	return posts, nil

	//posts, err := shared.MakeFacebookRequest(http.MethodGet, ctx.Auth.AccessToken, url, nil)
	//if err != nil {
	//	return nil, err
	//}
	//return posts, nil
}

func NewGetPostsAction() integration.Action {
	return &GetPostsAction{}
}
