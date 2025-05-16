package actions

import (
	"errors"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/mailchimp/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type addSubscriberToTagActionProps struct {
	ListID string `json:"list-id"`
	Tags   string `json:"tag-names"`
	Email  string `json:"email"`
}

type AddSubscriberToTagAction struct{}

func (a *AddSubscriberToTagAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		DisplayName:   "Add Subscriber To Tag",
		Description:   "Adds a subscriber to a specific tag in your email marketing platform, allowing you to easily manage and segment your audience.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: addSubscriberToTagDocs,
		SampleOutput: map[string]any{
			"success": "Subscriber added to tag!",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *AddSubscriberToTagAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("add_subscriber_to_tag", "Add Subscriber To Tag")

	form.TextField("list-id", "List ID").
		Placeholder("Enter a value for List ID.").
		Required(true).
		HelpText("The ID of the list to add the contact to.")

	form.TextField("tag-names", "Tag Names").
		Placeholder("Enter a value for Tag Names.").
		Required(true).
		HelpText("The tag names to add to the subscriber.")

	form.TextField("email", "Email").
		Placeholder("Enter a value for Email.").
		Required(true).
		HelpText("The email of the subscriber.")

	schema := form.Build()

	return schema
}

func (a *AddSubscriberToTagAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[addSubscriberToTagActionProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	dc, err := shared.GetMailChimpServerPrefix(token)
	if err != nil {
		return nil, fmt.Errorf("unable to get mailchimp server prefix: %w", err)
	}

	tags := shared.ProcessTagNamesInput(input.Tags)

	err = shared.ModifySubscriberTags(token, dc, input.ListID, input.Email, tags, "active")
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"status": "Tag added!",
	}, nil
}

func (a *AddSubscriberToTagAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewAddSubscriberToTagAction() sdk.Action {
	return &AddSubscriberToTagAction{}
}
