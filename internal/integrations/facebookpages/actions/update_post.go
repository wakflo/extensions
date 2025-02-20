package actions

import (
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/facebookpages/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"

	"github.com/wakflo/go-sdk/sdk"
)

type updatePostActionProps struct {
	PageID  string `json:"page_id"`
	PostID  string `json:"post_id"`
	Message string `json:"message"`
}

type UpdatePostAction struct{}

func (c *UpdatePostAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (c UpdatePostAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c UpdatePostAction) Name() string {
	return "Update a Post"
}

func (c UpdatePostAction) Description() string {
	return "update a post"
}

func (c UpdatePostAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &updatePostDocs,
	}
}

func (c UpdatePostAction) Icon() *string {
	return nil
}

func (c UpdatePostAction) SampleData() sdkcore.JSON {
	return nil
}

func (c UpdatePostAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"page_id": shared.GetFacebookPageInput("Select a page", "The page you want to get the post from", true),
		"post_id": shared.GetPagePostsInput("Select a post", "The post you want to update", true),
		"message": autoform.NewLongTextField().
			SetDisplayName("New Post").
			SetDescription("The text content of the post you want to edit").
			SetRequired(true).
			Build(),
	}
}

func (c UpdatePostAction) Auth() *sdk.Auth {
	return &sdk.Auth{
		Inherit: true,
	}
}

func (c UpdatePostAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updatePostActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	body := map[string]interface{}{
		"message": input.Message,
	}

	pageAccessToken, err := shared.GetPageAccessToken(ctx.Auth.AccessToken, input.PageID)
	if err != nil {
		return nil, err
	}

	url := "/" + input.PostID

	postResult, err := shared.PostActionFunc(pageAccessToken, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	return postResult, nil
}

func NewUpdatePostAction() sdk.Action {
	return &UpdatePostAction{}
}
