package survey_monkey

import (
	"errors"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createContactListProps struct {
	Name string `json:"name"`
}

type CreateContactListOperation struct {
	options *sdk.OperationInfo
}

func NewCreateContactListOperation() *CreateContactListOperation {
	return &CreateContactListOperation{
		options: &sdk.OperationInfo{
			Name:        "Create New Contact List",
			Description: "Create a new contact list",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"name": autoform.NewShortTextField().
					SetDisplayName("Name of contact list").
					SetDescription("Name of contact list").
					SetRequired(true).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *CreateContactListOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing calendly auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[createContactListProps](ctx)

	if input.Name == "" {
		return nil, errors.New("name is required")
	}

	contactList, _ := createContactList(accessToken, input.Name)

	return contactList, nil
}

func (c *CreateContactListOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateContactListOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
