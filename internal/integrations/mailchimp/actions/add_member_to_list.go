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

type addMemberToListActionProps struct {
	ListID    string `json:"list-id"`
	To        string `json:"to"`
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
	Email     string `json:"email"`
	Status    string `json:"status"`
}

type AddMemberToListAction struct{}

func (a *AddMemberToListAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		DisplayName:   "Add Member To List",
		Description:   "Adds a new member to an existing list in your workflow, allowing you to easily manage and track team members, stakeholders, or other relevant parties involved in the process. This integration action enables seamless addition of new members to lists, streamlining collaboration and communication within your workflow.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: addMemberToListDocs,
		SampleOutput: map[string]any{
			"success": "Contact Added!",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *AddMemberToListAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("add_member_to_list", "Add Member To List")

	form.TextField("list-id", "List ID").
		Placeholder("Enter a value for List ID.").
		Required(true).
		HelpText("The ID of the list to add the contact to.")

	form.TextField("first-name", "First Name").
		Placeholder("Enter a value for First Name.").
		Required(true).
		HelpText("The first name of the new contact.")

	form.TextField("last-name", "Last Name").
		Placeholder("Enter a value for Last Name.").
		Required(true).
		HelpText("The last name of the new contact.")

	form.TextField("email", "Email").
		Placeholder("Enter a value for Email.").
		Required(true).
		HelpText("The email of the new contact.")

	form.SelectField("status", "Status").
		Placeholder("Enter a value for Status.").
		Required(true).
		HelpText("The status of the new contact.").
		AddOptions(shared.MailchimpStatusType...)

	schema := form.Build()

	return schema
}

func (a *AddMemberToListAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[addMemberToListActionProps](ctx)
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
		return nil, errors.New(fmt.Sprintf("unable to get mailchimp server prefix: %v", err))
	}

	err = shared.AddContactToList(token, dc, input.ListID, input.Email, input.FirstName, input.Status, input.LastName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": "Contact Added!",
	}, err
}

func (a *AddMemberToListAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewAddMemberToListAction() sdk.Action {
	return &AddMemberToListAction{}
}
