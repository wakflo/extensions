package actions

import (
	"errors"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/mailchimp/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type updateSubscriberStatusActionProps struct {
	ListID string `json:"list-id"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type UpdateSubscriberStatusAction struct{}

func (a *UpdateSubscriberStatusAction) Name() string {
	return "Update Subscriber Status"
}

func (a *UpdateSubscriberStatusAction) Description() string {
	return "Updates the status of a subscriber in your application or database, allowing you to reflect changes in their subscription level, account information, or other relevant details."
}

func (a *UpdateSubscriberStatusAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *UpdateSubscriberStatusAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &updateSubscriberStatusDocs,
	}
}

func (a *UpdateSubscriberStatusAction) Icon() *string {
	return nil
}

func (a *UpdateSubscriberStatusAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"list-id": autoform.NewShortTextField().
			SetDisplayName(" List ID").
			SetDescription("").
			SetRequired(true).
			Build(),
		"email": autoform.NewShortTextField().
			SetDisplayName("Email").
			SetDescription("Email of the subscriber").
			SetRequired(true).
			Build(),
		"status": autoform.NewSelectField().
			SetDisplayName("Status").
			SetOptions(shared.MailchimpSubscriberStatus).
			SetRequired(true).
			Build(),
	}
}

func (a *UpdateSubscriberStatusAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateSubscriberStatusActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing mailchimp auth token")
	}
	accessToken := ctx.Auth.AccessToken

	dc, err := shared.GetMailChimpServerPrefix(accessToken)
	if err != nil {
		return nil, fmt.Errorf("unable to get mailchimp server prefix: %w", err)
	}

	err = shared.UpdateSubscriberStatus(accessToken, dc, input.ListID, input.Email, input.Status)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": "Subscriber status updated successfully!",
	}, err
}

func (a *UpdateSubscriberStatusAction) Auth() *sdk.Auth {
	return nil
}

func (a *UpdateSubscriberStatusAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *UpdateSubscriberStatusAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewUpdateSubscriberStatusAction() sdk.Action {
	return &UpdateSubscriberStatusAction{}
}
