package actions

import (
	"fmt"
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
		"status": autoform.NewSelectField().
			SetDisplayName("Contact Status").
			SetDescription("Status of the contact").
			SetOptions([]*sdkcore.AutoFormSchema{
				{Const: "ACTIVE", Title: "Active"},
				{Const: "INACTIVE", Title: "Inactive"},
				{Const: "UNSUBSCRIBED", Title: "Unsubscribed"},
			}).
			Build(),
	}
}

func (a *UpdateContactAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateContactActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	// Validate that contact ID is provided
	if input.ContactID == "" {
		return nil, fmt.Errorf("contact ID is required for updating a contact")
	}

	// Construct contact update data
	contactData := make(map[string]interface{})

	// Add optional fields conditionally
	if input.GivenName != "" {
		contactData["given_name"] = input.GivenName
	}
	if input.FamilyName != "" {
		contactData["family_name"] = input.FamilyName
	}
	if input.Email != "" {
		contactData["email"] = input.Email
	}

	// Handle company information
	if input.CompanyName != "" {
		contactData["company"] = map[string]interface{}{
			"name": input.CompanyName,
		}
	}

	// Handle job title
	if input.JobTitle != "" {
		contactData["job_title"] = input.JobTitle
	}

	// Handle phone number
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

	// Handle address
	if input.AddressLine1 != "" {
		addressType := "OTHER"
		if input.AddressType != "" {
			addressType = input.AddressType
		}

		address := map[string]string{
			"type":  addressType,
			"line1": input.AddressLine1,
		}

		// Add optional address fields
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

	// Handle opt-in reason
	if input.OptInReason != "" {
		contactData["opt_in_reason"] = input.OptInReason
	}

	// Handle contact status
	if input.Status != "" {
		contactData["status"] = input.Status
	}

	// Construct the endpoint URL with contact ID
	endpoint := fmt.Sprintf("/contacts/%s", input.ContactID)

	// Make the PUT request to update the contact
	updatedContact, err := shared.MakeKeapRequest(ctx.Auth.AccessToken, http.MethodPut, endpoint, contactData)
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
