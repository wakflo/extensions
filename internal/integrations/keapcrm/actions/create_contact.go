package actions

import (
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/keapcrm/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
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

func (a *CreateContactAction) Name() string {
	return "Create Contact"
}

func (a *CreateContactAction) Description() string {
	return "Create a new contact in Keap with specified details"
}

func (a *CreateContactAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateContactAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createContactDocs,
	}
}

func (c CreateContactAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"given_name": autoform.NewShortTextField().
			SetDisplayName("First Name").
			SetDescription("Contact's first name").
			SetRequired(true).Build(),
		"family_name": autoform.NewShortTextField().
			SetDisplayName("Last Name").
			SetDescription("Contact's last name").
			SetRequired(true).Build(),
		"email": autoform.NewShortTextField().
			SetDisplayName("Email").
			SetDescription("Contact's email address").
			SetRequired(true).Build(),
		"phone_number": autoform.NewShortTextField().
			SetDisplayName("Phone Number").
			SetDescription("Contact's phone number").
			SetRequired(false).Build(),
		"phone_type": autoform.NewSelectField().
			SetDisplayName("Phone Type").
			SetDescription("Type of phone number").
			SetOptions([]*sdkcore.AutoFormSchema{
				{Const: "MOBILE", Title: "Mobile"},
				{Const: "HOME", Title: "Home"},
				{Const: "WORK", Title: "Work"},
				{Const: "OTHER", Title: "Other"},
			}).
			Build(),
		"address_line1": autoform.NewShortTextField().
			SetDisplayName("Address Line 1").
			SetDescription("First line of contact's address").
			SetRequired(false).Build(),
		"address_line2": autoform.NewShortTextField().
			SetDisplayName("Address Line 2").
			SetDescription("Second line of contact's address").
			SetRequired(false).Build(),
		"locality": autoform.NewShortTextField().
			SetDisplayName("City").
			SetDescription("City/Locality of contact's address").
			SetRequired(false).Build(),
		"region": autoform.NewShortTextField().
			SetDisplayName("State/Region").
			SetDescription("State/Region of contact's address").
			SetRequired(false).Build(),
		"postal_code": autoform.NewShortTextField().
			SetDisplayName("Postal Code").
			SetDescription("Postal/ZIP code of contact's address").
			SetRequired(false).Build(),
		"country": autoform.NewShortTextField().
			SetDisplayName("Country").
			SetDescription("Country of contact's address").
			SetRequired(false).Build(),
		"address_type": autoform.NewSelectField().
			SetDisplayName("Address Type").
			SetDescription("Type of address.").
			SetOptions([]*sdkcore.AutoFormSchema{
				{Const: "BILLING", Title: "Billing"},
				{Const: "SHIPPING", Title: "Shipping"},
				{Const: "HOME", Title: "Home"},
				{Const: "WORK", Title: "Work"},
				{Const: "OTHER", Title: "Other"},
			}).
			Build(),
		"company_name": autoform.NewShortTextField().
			SetDisplayName("Company Name").
			SetDescription("Name of the contact's company").
			SetRequired(false).Build(),
		"job_title": autoform.NewShortTextField().
			SetDisplayName("Job Title").
			SetDescription("Contact's job title").
			SetRequired(false).Build(),
		"opt_in_reason": autoform.NewShortTextField().
			SetDisplayName("Opt-in Reason").
			SetDescription("Reason for the contact opting in to communications").
			SetRequired(false).Build(),
	}
}

func (a *CreateContactAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createContactActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	contactData := map[string]interface{}{
		"given_name":  input.GivenName,
		"family_name": input.FamilyName,
		"email":       input.Email,
	}

	if input.CompanyName != "" {
		contactData["company"] = map[string]interface{}{
			"name": input.CompanyName,
		}
	}

	if input.JobTitle != "" {
		contactData["job_title"] = input.JobTitle
	}

	if input.PhoneNumber != "" {
		phoneType := "OTHER"
		if input.PhoneType != "" {
			phoneType = input.PhoneType
		}

		contactData["phone_numbers"] = []map[string]string{
			{
				"type":   phoneType,
				"number": input.PhoneNumber,
			},
		}
	}

	if input.AddressLine1 != "" {
		addressType := "OTHER"
		if input.AddressType != "" {
			addressType = input.AddressType
		}

		address := map[string]string{
			"type":  addressType,
			"line1": input.AddressLine1,
		}

		if input.AddressLine2 != "" {
			address["line2"] = input.AddressLine2
		}
		if input.Locality != "" {
			address["locality"] = input.Locality
		}
		if input.Region != "" {
			address["region"] = input.Region
		}
		if input.PostalCode != "" {
			address["postal_code"] = input.PostalCode
		}
		if input.Country != "" {
			address["country"] = input.Country
		}

		contactData["addresses"] = []map[string]string{address}
	}

	// Add opt-in reason if provided
	if input.OptInReason != "" {
		contactData["opt_in_reason"] = input.OptInReason
	}

	// Create the contact
	endpoint := "/contacts"
	contact, err := shared.MakeKeapRequest(ctx.Auth.AccessToken, http.MethodPost, endpoint, contactData)
	if err != nil {
		return nil, err
	}

	return contact, nil
}

func (a *CreateContactAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateContactAction) SampleData() sdkcore.JSON {
	return map[string]any{
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
	}
}

func (a *CreateContactAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateContactAction() sdk.Action {
	return &CreateContactAction{}
}

func (a *CreateContactAction) Icon() *string {
	icon := "mdi:account-plus"
	return &icon
}
