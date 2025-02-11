package actions

import (
	"github.com/wakflo/go-sdk/autoform"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/facebookpages/shared"
	sdkcore "github.com/wakflo/go-sdk/core"

	"github.com/wakflo/go-sdk/integration"
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
	return "update a post."
}

func (c UpdatePostAction) Documentation() *integration.OperationDocumentation {
	return &integration.OperationDocumentation{
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
		"post_id": autoform.NewShortTextField().
			SetDisplayName("Page Post ID").
			SetDescription("Post ID").
			SetRequired(true).
			Build(),
		"message": autoform.NewShortTextField().
			SetDisplayName("Message").
			SetDescription("The text content of your post").
			SetRequired(true).
			Build(),
	}
}

func (c UpdatePostAction) Auth() *integration.Auth {
	return &integration.Auth{
		Inherit: true,
	}
}

func (c UpdatePostAction) Perform(ctx integration.PerformContext) (sdkcore.JSON, error) {
	input, err := integration.InputToTypeSafely[updatePostActionProps](ctx.BaseContext)
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

func NewUpdatePostAction() integration.Action {
	return &UpdatePostAction{}
}
