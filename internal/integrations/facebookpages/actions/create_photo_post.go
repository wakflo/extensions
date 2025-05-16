package actions

import (
	"fmt"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/facebookpages/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createPhotoPostActionProps struct {
	PageID  string `json:"page_id"`
	Url     string `json:"url"`
	Caption string `json:"caption"`
}

type CreatePhotoPostAction struct{}

// Metadata returns metadata about the action
func (c *CreatePhotoPostAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_page_photo",
		DisplayName:   "Create Page Photo",
		Description:   "Create a photo on a Facebook Page you manage",
		Type:          core.ActionTypeAction,
		Documentation: createPhotoPostDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"id":      "12345678901234567",
			"post_id": "12345678901234567_98765432109876543",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (c *CreatePhotoPostAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_page_photo", "Create Page Photo")

	shared.RegisterFacebookPageProps(form)

	form.TextField("url", "url").
		Placeholder("Photo").
		HelpText("A URL we can access for the photo").
		Required(true)

	form.TextareaField("caption", "caption").
		Placeholder("Caption").
		HelpText("The caption of the photo post").
		Required(false)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (c *CreatePhotoPostAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (c *CreatePhotoPostAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createPhotoPostActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	body := map[string]interface{}{
		"url": input.Url,
	}

	if input.Caption != "" {
		body["caption"] = input.Caption
	}

	pageAccessToken, err := shared.GetPageAccessToken(authCtx.Token.AccessToken, input.PageID)
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
