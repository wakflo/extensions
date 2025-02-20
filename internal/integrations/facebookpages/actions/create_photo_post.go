package actions

import (
	"fmt"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/facebookpages/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"

	"github.com/wakflo/go-sdk/sdk"
)

type createPhotoPostActionProps struct {
	PageID  string `json:"page_id"`
	Url     string `json:"url"`
	Caption string `json:"caption"`
}

type CreatePhotoPostAction struct{}

func (c *CreatePhotoPostAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (c CreatePhotoPostAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c CreatePhotoPostAction) Name() string {
	return "Create Page Photo"
}

func (c CreatePhotoPostAction) Description() string {
	return "Create a photo on a Facebook Page you manage"
}

func (c CreatePhotoPostAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createPhotoPostDocs,
	}
}

func (c CreatePhotoPostAction) Icon() *string {
	return nil
}

func (c CreatePhotoPostAction) SampleData() sdkcore.JSON {
	return nil
}

func (c CreatePhotoPostAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"page_id": shared.GetFacebookPageInput("Select Page", "The Page you want to post from", true),
		"url": autoform.NewShortTextField().
			SetDisplayName("Photo").
			SetDescription("A URL we can access for the photo").
			SetRequired(true).
			Build(),
		"caption": autoform.NewLongTextField().
			SetDisplayName("Caption").
			SetDescription("The caption of the photo post").
			SetRequired(false).
			Build(),
	}
}

func (c CreatePhotoPostAction) Auth() *sdk.Auth {
	return &sdk.Auth{
		Inherit: true,
	}
}

func (c CreatePhotoPostAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createPhotoPostActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}
	body := map[string]interface{}{
		"url": input.Url,
	}

	if input.Caption != "" {
		body["caption"] = input.Caption
	}

	pageAccessToken, err := shared.GetPageAccessToken(ctx.Auth.AccessToken, input.PageID)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("/%s/photos", input.PageID)
	photoResult, err := shared.PostActionFunc(pageAccessToken, http.MethodPost, endpoint, body)
	if err != nil {
		return nil, err
	}

	return photoResult, nil
}

func NewCreatePhotoPostAction() sdk.Action {
	return &CreatePhotoPostAction{}
}
