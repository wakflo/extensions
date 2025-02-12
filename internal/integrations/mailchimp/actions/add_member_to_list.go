package actions

import (
	"errors"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/mailchimp/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type addMemberToListActionProps struct {
	ListID    string `json:"list-id"`
	To        string `json:"to"`
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
	Email     string `json:"email"`
	Status    string `json:"status"`
}

type AddMemberToListAction struct{}

func (a *AddMemberToListAction) Name() string {
	return "Add Member To List"
}

func (a *AddMemberToListAction) Description() string {
	return "Adds a new member to an existing list in your workflow, allowing you to easily manage and track team members, stakeholders, or other relevant parties involved in the process. This integration action enables seamless addition of new members to lists, streamlining collaboration and communication within your workflow."
}

func (a *AddMemberToListAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *AddMemberToListAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &addMemberToListDocs,
	}
}

func (a *AddMemberToListAction) Icon() *string {
	return nil
}

func (a *AddMemberToListAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"list-id": autoform.NewShortTextField().
			SetDisplayName(" List ID").
			SetDescription("").
			SetRequired(true).
			Build(),
		"first-name": autoform.NewShortTextField().
			SetDisplayName(" First Name").
			SetDescription("First name of the new contact").
			SetRequired(true).
			Build(),
		"last-name": autoform.NewShortTextField().
			SetDisplayName(" Last Name").
			SetDescription("Last name of the new contact").
			SetRequired(true).
			Build(),
		"email": autoform.NewShortTextField().
			SetDisplayName("Email").
			SetDescription("Email of the new contact").
			SetRequired(true).
			Build(),
		"status": autoform.NewSelectField().
			SetDisplayName("Status").
			SetOptions(shared.MailchimpStatusType).
			SetRequired(true).
			Build(),
	}
}

func (a *AddMemberToListAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[addMemberToListActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing mailchimp auth token")
	}
	accessToken := ctx.Auth.AccessToken

	dc, err := shared.GetMailChimpServerPrefix(accessToken)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to get mailchimp server prefix: %v", err))
	}

	err = shared.AddContactToList(accessToken, dc, input.ListID, input.Email, input.FirstName, input.Status, input.LastName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": "Contact Added!",
	}, err
}

func (a *AddMemberToListAction) Auth() *sdk.Auth {
	return nil
}

func (a *AddMemberToListAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *AddMemberToListAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewAddMemberToListAction() sdk.Action {
	return &AddMemberToListAction{}
}
