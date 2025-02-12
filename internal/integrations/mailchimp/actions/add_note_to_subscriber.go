package actions

import (
	"errors"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/mailchimp/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type addNoteToSubscriberActionProps struct {
	ListID string `json:"list-id"`
	Email  string `json:"email"`
	Note   string `json:"note"`
}

type AddNoteToSubscriberAction struct{}

func (a *AddNoteToSubscriberAction) Name() string {
	return "Add Note To Subscriber"
}

func (a *AddNoteToSubscriberAction) Description() string {
	return "Adds a note to a subscriber's record, allowing you to store additional information and context about the subscriber. This integration action is useful for tracking important details or updates related to a specific subscriber, making it easier to manage and analyze their activity over time."
}

func (a *AddNoteToSubscriberAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *AddNoteToSubscriberAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &addNoteToSubscriberDocs,
	}
}

func (a *AddNoteToSubscriberAction) Icon() *string {
	return nil
}

func (a *AddNoteToSubscriberAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"list-id": autoform.NewShortTextField().
			SetDisplayName(" List ID").
			SetDescription("").
			SetRequired(true).
			Build(),
		"note": autoform.NewLongTextField().
			SetDisplayName(" Note").
			SetDescription("Note to add to the subscriber").
			SetRequired(true).
			Build(),
		"email": autoform.NewShortTextField().
			SetDisplayName("Email").
			SetDescription("Email of the subscriber").
			SetRequired(true).
			Build(),
	}
}

func (a *AddNoteToSubscriberAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[addNoteToSubscriberActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing mailchimp auth token")
	}
	accessToken := ctx.Auth.AccessToken

	dc, err := shared.GetMailChimpServerPrefix(accessToken)
	if err != nil {
		return nil, fmt.Errorf("unable to get mailchimp server prefix: %v", err)
	}

	err = shared.AddMemberNote(accessToken, dc, input.ListID, input.Email, input.Note)
	if err != nil {
		return nil, err
	}

	return sdk.JSON(map[string]interface{}{
		"result": "note added!",
	}), err
}

func (a *AddNoteToSubscriberAction) Auth() *sdk.Auth {
	return nil
}

func (a *AddNoteToSubscriberAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"result": "note added!",
	}
}

func (a *AddNoteToSubscriberAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewAddNoteToSubscriberAction() sdk.Action {
	return &AddNoteToSubscriberAction{}
}
