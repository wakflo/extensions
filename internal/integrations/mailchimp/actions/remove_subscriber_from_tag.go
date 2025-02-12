package actions

import (
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/mailchimp/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type removeSubscriberFromTagActionProps struct {
	ListID string `json:"list-id"`
	Tags   string `json:"tag-names"`
	Email  string `json:"email"`
}

type RemoveSubscriberFromTagAction struct{}

func (a *RemoveSubscriberFromTagAction) Name() string {
	return "Remove Subscriber From Tag"
}

func (a *RemoveSubscriberFromTagAction) Description() string {
	return "Remove Subscriber From Tag: This integration action removes a subscriber from a specific tag in your email marketing platform, ensuring that the individual is no longer associated with the designated group."
}

func (a *RemoveSubscriberFromTagAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *RemoveSubscriberFromTagAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &removeSubscriberFromTagDocs,
	}
}

func (a *RemoveSubscriberFromTagAction) Icon() *string {
	return nil
}

func (a *RemoveSubscriberFromTagAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"list-id": autoform.NewShortTextField().
			SetDisplayName(" List ID").
			SetDescription("").
			SetRequired(true).
			Build(),
		"tag-names": autoform.NewLongTextField().
			SetDisplayName(" Tag name").
			SetDescription("Tag name to remove from the subscriber").
			SetRequired(true).
			Build(),
		"email": autoform.NewShortTextField().
			SetDisplayName("Email").
			SetDescription("Email of the subscriber").
			SetRequired(true).
			Build(),
	}
}

func (a *RemoveSubscriberFromTagAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[removeSubscriberFromTagActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	accessToken := ctx.Auth.AccessToken

	dc, err := shared.GetMailChimpServerPrefix(accessToken)
	if err != nil {
		return nil, fmt.Errorf("unable to get mailchimp server prefix: %w", err)
	}

	tags := shared.ProcessTagNamesInput(input.Tags)

	err = shared.ModifySubscriberTags(accessToken, dc, input.ListID, input.Email, tags, "inactive")
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"status": "Tag Removed!",
	}, nil
}

func (a *RemoveSubscriberFromTagAction) Auth() *sdk.Auth {
	return nil
}

func (a *RemoveSubscriberFromTagAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *RemoveSubscriberFromTagAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewRemoveSubscriberFromTagAction() sdk.Action {
	return &RemoveSubscriberFromTagAction{}
}
