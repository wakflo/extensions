package actions

import (
	"errors"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/mailchimp/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type addSubscriberToTagActionProps struct {
	ListID string `json:"list-id"`
	Tags   string `json:"tag-names"`
	Email  string `json:"email"`
}

type AddSubscriberToTagAction struct{}

func (a *AddSubscriberToTagAction) Name() string {
	return "Add Subscriber To Tag"
}

func (a *AddSubscriberToTagAction) Description() string {
	return "Add Subscriber to Tag: Automatically adds one or more subscribers to a specific tag in your email marketing platform, allowing you to easily manage and segment your audience."
}

func (a *AddSubscriberToTagAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *AddSubscriberToTagAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &addSubscriberToTagDocs,
	}
}

func (a *AddSubscriberToTagAction) Icon() *string {
	return nil
}

func (a *AddSubscriberToTagAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"list-id": autoform.NewShortTextField().
			SetDisplayName(" List ID").
			SetDescription("").
			SetRequired(true).
			Build(),
		"tag-names": autoform.NewLongTextField().
			SetDisplayName(" Tag Name").
			SetDescription("Tag name to add to the subscriber").
			SetRequired(true).
			Build(),
		"email": autoform.NewShortTextField().
			SetDisplayName("Email").
			SetDescription("Email of the subscriber").
			SetRequired(true).
			Build(),
	}
}

func (a *AddSubscriberToTagAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[addSubscriberToTagActionProps](ctx.BaseContext)
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

	tags := shared.ProcessTagNamesInput(input.Tags)

	err = shared.ModifySubscriberTags(accessToken, dc, input.ListID, input.Email, tags, "active")
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"status": "Tag added!",
	}, nil
}

func (a *AddSubscriberToTagAction) Auth() *sdk.Auth {
	return nil
}

func (a *AddSubscriberToTagAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"status": "Tag added!",
	}
}

func (a *AddSubscriberToTagAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewAddSubscriberToTagAction() sdk.Action {
	return &AddSubscriberToTagAction{}
}
