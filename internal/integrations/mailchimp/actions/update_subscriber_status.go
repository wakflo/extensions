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

type updateSubscriberStatusActionProps struct {
	ListID string `json:"list-id"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type UpdateSubscriberStatusAction struct{}

// func (a *UpdateSubscriberStatusAction) Name() string {
// 	return "Update Subscriber Status"
// }

func (a *UpdateSubscriberStatusAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		DisplayName:   "Update Subscriber Status",
		Description:   "Update Subscriber Status: This integration action updates the status of a subscriber in your application or database, allowing you to reflect changes in their subscription level, account information, or other relevant details.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: updateSubscriberStatusDocs,
		SampleOutput: map[string]any{
			"success": "Subscriber status updated!",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *UpdateSubscriberStatusAction) Properties() *smartform.FormSchema {

	form := smartform.NewForm("update_subscriber_status", "Update Subscriber Status")

	form.TextField("list-id", "List ID").
		Placeholder("Enter a value for List ID.").
		Required(true).
		HelpText("The ID of the list to update the contact status for.")

	form.TextField("email", "Email").
		Placeholder("Enter a value for Email.").
		Required(true).
		HelpText("The email of the subscriber.")

	form.SelectField("status", "Status").
		Placeholder("Enter a value for Status.").
		Required(true).
		AddOptions(shared.MailchimpSubscriberStatus...).
		HelpText("The status to update the subscriber to.")

	schema := form.Build()

	return schema
}

func (a *UpdateSubscriberStatusAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateSubscriberStatusActionProps](ctx)
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

	err = shared.UpdateSubscriberStatus(token, dc, input.ListID, input.Email, input.Status)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": "Subscriber status updated successfully!",
	}, err
}

func (a *UpdateSubscriberStatusAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewUpdateSubscriberStatusAction() sdk.Action {
	return &UpdateSubscriberStatusAction{}
}
