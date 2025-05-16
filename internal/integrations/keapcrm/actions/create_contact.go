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

type createContactActionProps struct {
	GivenName    string `json:"given_name"`
	FamilyName   string `json:"family_name"`
	Email        string `json:"email"`
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
}

type CreateContactAction struct{}

func (a *CreateContactAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		DisplayName:   "Create Contact",
		Description:   "Create a new contact in Keap with specified details",
		Type:          sdkcore.ActionTypeAction,
		Documentation: createContactDocs,
		SampleOutput: map[string]any{
			"id":          "12345",
			"given_name":  "John",
			"family_name": "Doe",
			"email":       "john.doe@example.com",
			"phone_numbers": []map[string]string{
				{
					"type":   "WORK",
					"number": "+1 555-123-4567",
				},
			},
			"addresses": []map[string]string{
				{
					"type":        "BILLING",
					"line1":       "123 Main St",
					"line2":       "Suite 100",
					"locality":    "Anytown",
					"region":      "CA",
					"postal_code": "12345",
					"country":     "USA",
				},
			},
			"company": map[string]string{
				"name": "Acme Inc",
			},
			"job_title": "CEO",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (c CreateContactAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_contact", "Create Contact")

	form.TextField("given_name", "First Name").
		Placeholder("Enter a value for First Name.").
		Required(true).
		HelpText("Contact's first name.")

	form.TextField("family_name", "Last Name").
		Placeholder("Enter a value for Last Name.").
		Required(true).
		HelpText("Contact's last name.")

	form.TextField("email", "Work Email").
		Placeholder("Enter a value for Work Email.").
		Required(true).
		HelpText("Contact's email address.")

	form.TextField("phone_number", "Phone Number").
		Placeholder("Enter a value for Phone Number.").
		Required(false).
		HelpText("Contact's phone number.")

	schema := form.Build()

	return schema
}

func (a *CreateContactAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createContactActionProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	contactData := map[string]interface{}{
		"given_name":  input.GivenName,
		"family_name": input.FamilyName,
		"email_addresses": []map[string]string{
			{
				"email": input.Email,
				"field": "EMAIL1",
			},
		},
	}

	if input.PhoneNumber != "" {
		contactData["phone_numbers"] = []map[string]interface{}{
			{
				"number": input.PhoneNumber,
				"field":  "PHONE1",
			},
		}
	}

	endpoint := "/contacts"
	contact, err := shared.MakeKeapRequest(token, http.MethodPost, endpoint, contactData)
	if err != nil {
		return nil, err
	}

	return contact, nil
}

func (a *CreateContactAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewCreateContactAction() sdk.Action {
	return &CreateContactAction{}
}
