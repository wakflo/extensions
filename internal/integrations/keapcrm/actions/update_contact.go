package actions

import (
	"errors"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/keapcrm/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type updateContactActionProps struct {
	ContactID    string `json:"contact_id"`
	GivenName    string `json:"given_name,omitempty"`
	FamilyName   string `json:"family_name,omitempty"`
	Email        string `json:"email,omitempty"`
	PhoneNumber  string `json:"phone_number,omitempty"`
	PhoneType    string `json:"phone_type,omitempty"`
	AddressLine1 string `json:"address_line1,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty"`
	Locality     string `json:"locality,omitempty"`
	Region       string `json:"region,omitempty"`
	PostalCode   string `json:"postal_code,omitempty"`
	Country      string `json:"country,omitempty"`
	AddressType  string `json:"address_type,omitempty"`
	CompanyName  string `json:"company_name,omitempty"`
	JobTitle     string `json:"job_title,omitempty"`
	OptInReason  string `json:"opt_in_reason,omitempty"`
	Status       string `json:"status,omitempty"`
}

type UpdateContactAction struct{}

func (a *UpdateContactAction) Name() string {
	return "Update Contact"
}

func (a *UpdateContactAction) Description() string {
	return "Update an existing contact in Keap with specified details"
}

func (a *UpdateContactAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *UpdateContactAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &updateContactDocs,
	}
}

func (a *UpdateContactAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"contact_id": autoform.NewShortTextField().
			SetDisplayName("Contact ID").
			SetDescription("Unique identifier of the contact to update").
			SetRequired(true).Build(),
		"given_name": autoform.NewShortTextField().
			SetDisplayName("First Name").
			SetDescription("Contact's first name").
			SetRequired(false).Build(),
		"family_name": autoform.NewShortTextField().
			SetDisplayName("Last Name").
			SetDescription("Contact's last name").
			SetRequired(false).Build(),
		"email": autoform.NewShortTextField().
			SetDisplayName("Email").
			SetDescription("Contact's email address").
			SetRequired(false).Build(),
		"phone_number": autoform.NewShortTextField().
			SetDisplayName("Phone Number").
			SetDescription("Contact's phone number").
			SetRequired(false).Build(),
	}
}

func (a *UpdateContactAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateContactActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if input.ContactID == "" {
		return nil, errors.New("contact ID is required for updating a contact")
	}

	contactData := make(map[string]interface{})

	if input.GivenName != "" {
		contactData["given_name"] = input.GivenName
	}
	if input.FamilyName != "" {
		contactData["family_name"] = input.FamilyName
	}
	if input.Email != "" {
		contactData["email_addresses"] = []map[string]interface{}{
			{
				"email": input.Email,
				"field": "EMAIL1",
			},
		}
	}

	if input.PhoneNumber != "" {
		contactData["phone_numbers"] = []map[string]interface{}{
			{
				"field":  "PHONE1",
				"number": input.PhoneNumber,
			},
		}
	}

	endpoint := "/contacts/" + input.ContactID

	updatedContact, err := shared.MakeKeapRequest(ctx.Auth.AccessToken, http.MethodPatch, endpoint, contactData)
	if err != nil {
		return nil, err
	}

	return updatedContact, nil
}

func (a *UpdateContactAction) Auth() *sdk.Auth {
	return nil
}

func (a *UpdateContactAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":          "12345",
		"given_name":  "John",
		"family_name": "Doe",
		"email":       "john.doe@example.com",
		"status":      "ACTIVE",
	}
}

func (a *UpdateContactAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewUpdateContactAction() sdk.Action {
	return &UpdateContactAction{}
}

func (a *UpdateContactAction) Icon() *string {
	icon := "mdi:account-edit"
	return &icon
}
