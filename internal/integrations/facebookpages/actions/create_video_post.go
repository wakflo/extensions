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

type createVideoPostActionProps struct {
	PageID      string `json:"page_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Video       string `json:"video"`
}

type CreateVideoPostAction struct{}

// Metadata returns metadata about the action
func (c *CreateVideoPostAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_video_post",
		DisplayName:   "Create Video Post",
		Description:   "Create a video post on a Facebook Page you manage",
		Type:          core.ActionTypeAction,
		Documentation: createVideoPostDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"id":      "12345678901234567",
			"success": true,
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (c *CreateVideoPostAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_video_post", "Create Video Post")

	shared.RegisterFacebookPageProps(form)

	form.TextField("title", "title").
		Placeholder("Title").
		HelpText("Title of the video post").
		Required(false)

	form.TextField("video", "video").
		Placeholder("Video Url").
		HelpText("A URL we can access for the video (Limit: 1GB or 20 minutes)").
		Required(true)

	form.TextareaField("description", "description").
		Placeholder("Description").
		HelpText("Description of the video post").
		Required(false)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (c *CreateVideoPostAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (c *CreateVideoPostAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createVideoPostActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
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
	pageAccessToken, err := shared.GetPageAccessToken(authCtx.Token.AccessToken, input.PageID)
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
