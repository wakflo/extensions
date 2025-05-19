package actions

import (
	"errors"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/keapcrm/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
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

func (a *UpdateContactAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "update_contact",
		DisplayName:   "Update Contact",
		Description:   "Update an existing contact in Keap with specified details",
		Type:          sdkcore.ActionTypeAction,
		Documentation: updateContactDocs,
		SampleOutput: map[string]any{
			"id":          "12345",
			"given_name":  "John",
			"family_name": "Doe",
			"email":       "john.doe@example.com",
			"status":      "ACTIVE",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *UpdateContactAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("update_contact", "Update Contact")

	form.TextField("contact_id", "Contact ID").
		Placeholder("contact ID").
		Required(true).
		HelpText("Unique identifier of the contact to update")

	form.TextField("given_name", "First Name").
		Required(false).
		Placeholder("First Name").
		HelpText("Contact's first name")

	form.TextField("family_name", "Last Name").
		Required(false).
		Placeholder("Last Name").
		HelpText("Contact's last name")

	form.TextField("email", "Email").
		Required(false).
		Placeholder("Email").
		HelpText("Contact's email address")

	form.TextField("phone_number", "Phone Number").
		Required(false).
		Placeholder("Phone Number").
		HelpText("Contact's phone number")

	schema := form.Build()

	return schema
}

func (a *UpdateContactAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateContactActionProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

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

	updatedContact, err := shared.MakeKeapRequest(token, http.MethodPatch, endpoint, contactData)
	if err != nil {
		return nil, err
	}

	return updatedContact, nil
}

func (a *UpdateContactAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewUpdateContactAction() sdk.Action {
	return &UpdateContactAction{}
}
