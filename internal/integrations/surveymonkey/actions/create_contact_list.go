package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/surveymonkey/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type createContactListProps struct {
	Name string `json:"name"`
}

type CreateContactListAction struct{}

func (a *CreateContactListAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_contact_list",
		DisplayName:   "Create Contact List",
		Description:   "Create a new contact list",
		Type:          sdkcore.ActionTypeAction,
		Documentation: createContactListDocs,
		SampleOutput: map[string]any{
			"message": "hello world",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *CreateContactListAction) Properties() *smartform.FormSchema {

	form := smartform.NewForm("create_contact_list", "Create Contact List")

	form.TextField("name", "Contact List Name").
		Placeholder("Enter a value for Contact List Name.").
		Required(true).
		HelpText("The name of the contact list")

	schema := form.Build()

	return schema
}

func (a *CreateContactListAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createContactListProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	if input.Name == "" {
		return nil, errors.New("name is required")
	}

	resData, err := shared.CreateContactList(token, input.Name)
	if err != nil {
		return nil, err
	}

	return resData, nil
}

func (a *CreateContactListAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewCreateContactListAction() sdk.Action {
	return &CreateContactListAction{}
}
