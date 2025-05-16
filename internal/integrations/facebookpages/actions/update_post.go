package actions

import (
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/facebookpages/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type updatePostActionProps struct {
	PageID  string `json:"page_id"`
	PostID  string `json:"post_id"`
	Message string `json:"message"`
}

type UpdatePostAction struct{}

// Metadata returns metadata about the action
func (c *UpdatePostAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "update_post",
		DisplayName:   "Update a Post",
		Description:   "Update a post on a Facebook Page",
		Type:          core.ActionTypeAction,
		Documentation: updatePostDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"success": true,
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (c *UpdatePostAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("update_post", "Update a Post")

	shared.RegisterFacebookPageProps(form)

	shared.RegisterPagePostsProps(form)

	form.TextareaField("message", "message").
		Placeholder("New Content").
		HelpText("The text content of the post you want to edit").
		Required(true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (c *UpdatePostAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (c *UpdatePostAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[updatePostActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	body := map[string]interface{}{
		"message": input.Message,
	}

	pageAccessToken, err := shared.GetPageAccessToken(authCtx.AccessToken, input.PageID)
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
