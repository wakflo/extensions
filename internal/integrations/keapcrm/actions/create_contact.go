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
			SetDisplayName("Work Email").
			SetDescription("Contact's email address").
			SetRequired(true).Build(),
		"phone_number": autoform.NewShortTextField().
			SetDisplayName("Phone Number").
			SetDescription("Contact's phone number").
			SetRequired(false).Build(),
	}
}

func (a *CreateContactAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createContactActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

contactData := map[string]interface{}{
    "given_name":     input.GivenName,
    "family_name":    input.FamilyName,
    "email_addresses": []map[string]string{
        {
            "email":    input.Email,
            "field":    "EMAIL1",
        },
    },
}

	
if input.PhoneNumber != "" {

    contactData["phone_numbers"] = []map[string]interface{}{
        {
            "number": input.PhoneNumber,
            "field": "PHONE1",
        },
    }
}

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
