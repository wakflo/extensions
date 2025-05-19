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

type removeSubscriberFromTagActionProps struct {
	ListID string `json:"list-id"`
	Tags   string `json:"tag-names"`
	Email  string `json:"email"`
}

type RemoveSubscriberFromTagAction struct{}

func (a *RemoveSubscriberFromTagAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "remove_subscriber_from_tag",
		DisplayName:   "Remove Subscriber From Tag",
		Description:   "Remove Subscriber From Tag: This integration action removes a subscriber from a specific tag in your email marketing platform, ensuring that the individual is no longer associated with the designated group.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: removeSubscriberFromTagDocs,
		SampleOutput: map[string]any{
			"success": "Subscriber removed from tag!",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *RemoveSubscriberFromTagAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("remove_subscriber_from_tag", "Remove Subscriber From Tag")

	form.TextField("list-id", "List ID").
		Placeholder("Enter a value for List ID.").
		Required(true).
		HelpText("The ID of the list to remove the contact from.")

	form.TextField("tag-names", "Tag Names").
		Placeholder("Enter a value for Tag Names.").
		Required(true).
		HelpText("The tag names to remove from the subscriber.")

	form.TextField("email", "Email").
		Placeholder("Enter a value for Email.").
		Required(true).
		HelpText("The email of the subscriber.")

	schema := form.Build()

	return schema
}

func (a *RemoveSubscriberFromTagAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[removeSubscriberFromTagActionProps](ctx)
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

	err = shared.ModifySubscriberTags(token, dc, input.ListID, input.Email, tags, "inactive")
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"status": "Tag Removed!",
	}, nil
}

func (a *RemoveSubscriberFromTagAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewRemoveSubscriberFromTagAction() sdk.Action {
	return &RemoveSubscriberFromTagAction{}
}
