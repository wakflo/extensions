package actions

import (
	"fmt"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/facebookpages/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"

	"github.com/wakflo/go-sdk/sdk"
)

type createVideoPostActionProps struct {
	PageID      string `json:"page_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Video       string `json:"video"`
}

type CreateVideoPostAction struct{}

func (c *CreateVideoPostAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (c CreateVideoPostAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c CreateVideoPostAction) Name() string {
	return "Create Video Post"
}

func (c CreateVideoPostAction) Description() string {
	return "Create a video post on a Facebook Page you manage"
}

func (c CreateVideoPostAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createVideoPostDocs,
	}
}

func (c CreateVideoPostAction) Icon() *string {
	return nil
}

func (c CreateVideoPostAction) SampleData() sdkcore.JSON {
	return nil
}

func (c CreateVideoPostAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"page_id": shared.GetFacebookPageInput("Select a page", "The page you want to get the post from", true),
		"title": autoform.NewShortTextField().
			SetDisplayName("Title").
			SetDescription("Title of the video post").
			SetRequired(false).
			Build(),
		"video": autoform.NewShortTextField().
			SetDisplayName("Video Url").
			SetDescription("A URL we can access for the video (Limit: 1GB or 20 minutes)").
			SetRequired(true).
			Build(),
		"description": autoform.NewLongTextField().
			SetDisplayName("Title").
			SetDescription("Title of the video post").
			SetRequired(false).
			Build(),
	}
}

func (c CreateVideoPostAction) Auth() *sdk.Auth {
	return &sdk.Auth{
		Inherit: true,
	}
}

func (c CreateVideoPostAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createVideoPostActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}
	body := map[string]interface{}{
		"file_url": input.Video,
	}

	if input.Title != "" {
		body["title"] = input.Title
	}

	if input.Description != "" {
		body["description"] = input.Description
	}

	endpoint := fmt.Sprintf("/%s/videos", input.PageID)
	pageAccessToken, err := shared.GetPageAccessToken(ctx.Auth.AccessToken, input.PageID)
	if err != nil {
		return nil, err
	}
	videoResult, err := shared.PostActionFunc(pageAccessToken, http.MethodPost, endpoint, body)
	if err != nil {
		return nil, err
	}
	return videoResult, nil
}

func NewCreateVideoPostAction() sdk.Action {
	return &CreateVideoPostAction{}
}
