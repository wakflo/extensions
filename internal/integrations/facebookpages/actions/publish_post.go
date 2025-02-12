package actions

import (
	"fmt"
	"github.com/wakflo/extensions/internal/integrations/facebookpages/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"net/http"
	"strings"
	"time"

	"github.com/wakflo/go-sdk/integration"
)

type publishPostActionProps struct {
	PageID               string `json:"page_id"`
	Message              string `json:"message"`
	Link                 string `json:"link"`
	Published            *bool  `json:"published"`
	ScheduledPublishTime string `json:"scheduled_publish_time"`
}

type PublishPostAction struct{}

func (c *PublishPostAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (c PublishPostAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c PublishPostAction) Name() string {
	return "Publish Post to Facebook"
}

func (c PublishPostAction) Description() string {
	return "Publish a new post to your Facebook Page"
}

func (c PublishPostAction) Documentation() *integration.OperationDocumentation {
	return &integration.OperationDocumentation{
		Documentation: &publishPostDocs,
	}
}

func (c PublishPostAction) Icon() *string {
	return nil
}

func (c PublishPostAction) SampleData() sdkcore.JSON {
	return nil
}

func (c PublishPostAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"page_id": shared.GetFacebookPageInput("Select A Page", "The page you want to to publish post from", true),
		"message": autoform.NewLongTextField().
			SetDisplayName("Message").
			SetDescription("The text content of your post").
			SetRequired(true).
			Build(),
		"link": autoform.NewShortTextField().
			SetDisplayName("Link").
			SetDescription("URL to attach to the post").
			SetRequired(false).
			Build(),
		"published": autoform.NewBooleanField().
			SetDisplayName("Publish Immediately").
			SetDescription("Set to false to schedule the post for later").
			SetDefaultValue(true).
			Build(),
		"scheduled_publish_time": autoform.NewDateTimeField().
			SetDisplayName("Scheduled Publish Time").
			SetDescription("When to publish the post. Can be a UNIX timestamp, ISO 8601 timestamp, or relative time (e.g. '+2 weeks'). Required if not publishing immediately.").
			SetRequired(false).
			Build(),
	}
}

func (c PublishPostAction) Auth() *integration.Auth {
	return &integration.Auth{
		Inherit: true,
	}
}

func (c PublishPostAction) Perform(ctx integration.PerformContext) (sdkcore.JSON, error) {
	input, err := integration.InputToTypeSafely[publishPostActionProps](ctx.BaseContext)
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

	pageAccessToken, err := shared.GetPageAccessToken(ctx.Auth.AccessToken, input.PageID)
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

func NewPublishPostAction() integration.Action {
	return &PublishPostAction{}
}
