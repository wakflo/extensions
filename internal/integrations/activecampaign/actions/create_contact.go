package actions

import (
	"encoding/json"
	"errors"

	"github.com/wakflo/extensions/internal/integrations/activecampaign/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createContactActionProps struct {
	Email     string `json:"email"`
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
	Phone     string `json:"phone"`
}

type CreateContactAction struct{}

func (a *CreateContactAction) Name() string {
	return "Create Contact"
}

func (a *CreateContactAction) Description() string {
	return "Create a new contact in your ActiveCampaign account with customizable fields."
}

func (a *CreateContactAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateContactAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createContactDocs,
	}
}

func (a *CreateContactAction) Icon() *string {
	icon := "mdi:account-plus"
	return &icon
}

func (a *CreateContactAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"email": autoform.NewShortTextField().
			SetDisplayName("Email").
			SetDescription("Email address of the contact").
			SetRequired(true).
			Build(),
		"first-name": autoform.NewShortTextField().
			SetDisplayName("First Name").
			SetDescription("First name of the contact").
			SetRequired(false).
			Build(),
		"last-name": autoform.NewShortTextField().
			SetDisplayName("Last Name").
			SetDescription("Last name of the contact").
			SetRequired(false).
			Build(),
		"phone": autoform.NewShortTextField().
			SetDisplayName("Phone").
			SetDescription("Phone number of the contact").
			SetRequired(false).
			Build(),
		// "list-ids": autoform.NewShortTextField().
		// 	SetDisplayName("List IDs").
		// 	SetDescription("Comma-separated list of list IDs to add the contact to").
		// 	SetRequired(false).
		// 	Build(),

	}
}

func (a *CreateContactAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createContactActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if input.Email == "" {
		return nil, errors.New("email is required")
	}

	contactData := map[string]interface{}{
		"contact": map[string]interface{}{
			"email": input.Email,
		},
	}

	contact := contactData["contact"].(map[string]interface{})

	if input.FirstName != "" {
		contact["firstName"] = input.FirstName
	}

	if input.LastName != "" {
		contact["lastName"] = input.LastName
	}

	if input.Phone != "" {
		contact["phone"] = input.Phone
	}

	payload, err := json.Marshal(contactData)
	if err != nil {
		return nil, err
	}

	response, err := shared.PostActiveCampaignClient(
		ctx.Auth.Extra["api_url"],
		ctx.Auth.Extra["api_key"],
		"contacts",
		payload,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *CreateContactAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateContactAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":        "123",
		"email":     "new@example.com",
		"firstName": "New",
		"lastName":  "User",
		"phone":     "+1234567890",
		"cdate":     "2023-03-15T15:30:00-05:00",
		"udate":     "2023-03-15T15:30:00-05:00",
	}
}

func (a *CreateContactAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateContactAction() sdk.Action {
	return &CreateContactAction{}
}
