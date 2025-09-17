package actions

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/facebookpages/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type publishPostActionProps struct {
	PageID               string `json:"page_id"`
	Message              string `json:"message"`
	Link                 string `json:"link"`
	Published            *bool  `json:"published"`
	ScheduledPublishTime string `json:"scheduled_publish_time"`
}

type PublishPostAction struct{}

// Metadata returns metadata about the action
func (c *PublishPostAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "publish_post_to_facebook",
		DisplayName:   "Publish Post to Facebook",
		Description:   "Publish a new post to your Facebook Page",
		Type:          core.ActionTypeAction,
		Documentation: publishPostDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"id":      "12345678901234567_98765432109876543",
			"success": true,
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (c *PublishPostAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("publish_post_to_facebook", "Publish Post to Facebook")

	shared.RegisterFacebookPageProps(form)

	form.TextareaField("message", "message").
		Placeholder("Message").
		HelpText("The text content of your post").
		Required(true)

	form.TextField("link", "link").
		Placeholder("Link").
		HelpText("URL to attach to the post").
		Required(false)

	form.CheckboxField("published", "published").
		Placeholder("Publish Immediately").
		HelpText("Set to false to schedule the post for later").
		DefaultValue(true).
		Required(false)

	form.DateField("scheduled_publish_time", "scheduled_publish_time").
		Placeholder("Scheduled Publish Time").
		HelpText("When to publish the post. Can be a UNIX timestamp, ISO 8601 timestamp, or relative time (e.g. '+2 weeks'). Required if not publishing immediately.").
		Required(false)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (c *PublishPostAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (c *PublishPostAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[publishPostActionProps](ctx)
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

	if input.Link != "" {
		body["link"] = input.Link
	}

	if input.Published != nil {
		body["published"] = *input.Published

		if !*input.Published && input.ScheduledPublishTime == "" {
			return nil, fmt.Errorf("scheduled_publish_time is required when published is false")
		}

		if input.ScheduledPublishTime != "" {
			cleanTimeString := strings.Split(input.ScheduledPublishTime, "[")[0]

			parsedTime, errs := time.Parse(time.RFC3339, cleanTimeString)
			if errs != nil {
				return nil, fmt.Errorf("error parsing scheduled_publish_time: %v", err)
			}

			scheduledPublishTime := parsedTime.Format("2006-01-02T15:04:05-07:00")
			body["scheduled_publish_time"] = scheduledPublishTime
		}
	}

	pageAccessToken, err := shared.GetPageAccessToken(authCtx.Token.AccessToken, input.PageID)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("/%s/feed", input.PageID)

	posts, err := shared.PostActionFunc(pageAccessToken, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func NewPublishPostAction() sdk.Action {
	return &PublishPostAction{}
}
