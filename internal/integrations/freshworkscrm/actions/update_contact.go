package actions

import (
	"errors"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/freshworkscrm/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type updateContactActionProps struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	MobileNumber  string `json:"mobile_number"`
	Email         string `json:"email"`
	ContactViewID string `json:"contact_view_id"`
	ContactID     string `json:"contact_id"`
}

type UpdateContactAction struct{}

func (a *UpdateContactAction) Name() string {
	return "Update a Contact"
}

func (a *UpdateContactAction) Description() string {
	return "Update a contact in Freshworks CRM with specified information."
}

func (a *UpdateContactAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *UpdateContactAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &updateContactDocs,
	}
}

func (a *UpdateContactAction) Icon() *string {
	return nil
}

func (a *UpdateContactAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"first_name": autoform.NewShortTextField().
			SetDisplayName("First Name").
			SetDescription("Contact's first name").
			SetRequired(false).
			Build(),
		"last_name": autoform.NewShortTextField().
			SetDisplayName("Last Name").
			SetDescription("Contact's last name").
			SetRequired(false).
			Build(),
		"email": autoform.NewShortTextField().
			SetDisplayName("Email").
			SetDescription("Contact's email").
			SetRequired(false).
			Build(),
		"mobile_number": autoform.NewShortTextField().
			SetDisplayName("Mobile Number").
			SetDescription("Contact's mobile number").
			SetRequired(false).
			Build(),
		"contact_view_id": shared.GetContactViewInput(),
		"contact_id":      shared.GetContactsInput(),
	}
}

func (a *UpdateContactAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" || ctx.Auth.Extra["domain"] == "" {
		return nil, errors.New("missing freshworks auth parameters")
	}

	input := sdk.InputToType[updateContactActionProps](ctx.BaseContext)

	domain := ctx.Auth.Extra["domain"]
	freshworksDomain := "https://" + domain + ".myfreshworks.com"

	contact := make(map[string]interface{})

	shared.UpdateField(contact, "first_name", input.FirstName)
	shared.UpdateField(contact, "last_name", input.LastName)
	shared.UpdateField(contact, "email", input.Email)
	shared.UpdateField(contact, "mobile_number", input.MobileNumber)

	contactData := map[string]interface{}{
		"contact": contact,
	}

	response, err := shared.UpdateContact(freshworksDomain, ctx.Auth.Extra["api-key"], input.ContactID, contactData)
	if err != nil {
		return nil, fmt.Errorf("error updating contact:  %v", err)
	}

	return response, nil
}

func (a *UpdateContactAction) Auth() *sdk.Auth {
	return nil
}

func (a *UpdateContactAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"contact": map[string]any{
			"id":         "12345",
			"first_name": "John",
			"last_name":  "Doe",
			"email":      "john.doe@example.com",
			"created_at": "2023-01-01T12:00:00Z",
			"updated_at": "2023-01-01T12:00:00Z",
		},
	}
}

func (a *UpdateContactAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewUpdateContactAction() sdk.Action {
	return &UpdateContactAction{}
}
