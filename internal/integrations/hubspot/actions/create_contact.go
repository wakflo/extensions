package actions

import (
	"encoding/json"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/hubspot/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createContactActionProps struct {
	FirstName  string `json:"firstname"`
	LastName   string `json:"lastname"`
	Email      string `json:"email"`
	Zip        string `json:"zip"`
	Phone      string `json:"phone"`
	Company    string `json:"company"`
	JobTitle   string `json:"jobtitle"`
	Website    string `json:"website"`
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	ZipCode    string `json:"zipcode"`
	Country    string `json:"country"`
	Properties string `json:"properties"`
}

type CreateContactAction struct{}

func (a *CreateContactAction) Name() string {
	return "Create Contact"
}

func (a *CreateContactAction) Description() string {
	return "Create a new contact in HubSpot with the specified information."
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
	return nil
}

func (a *CreateContactAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"email": autoform.NewShortTextField().
			SetDisplayName("Email").
			SetDescription("Email address of the contact").
			SetRequired(true).
			Build(),
		"firstname": autoform.NewShortTextField().
			SetDisplayName("First Name").
			SetDescription("First name of the contact").
			SetRequired(false).
			Build(),
		"lastname": autoform.NewShortTextField().
			SetDisplayName("Last Name").
			SetDescription("Last name of the contact").
			SetRequired(false).
			Build(),
		"phone": autoform.NewShortTextField().
			SetDisplayName("Phone").
			SetDescription("Phone number of the contact").
			SetRequired(false).
			Build(),
		"company": autoform.NewShortTextField().
			SetDisplayName("Company").
			SetDescription("Company name the contact works for").
			SetRequired(false).
			Build(),
		"jobtitle": autoform.NewShortTextField().
			SetDisplayName("Job Title").
			SetDescription("Job title of the contact").
			SetRequired(false).
			Build(),
		"website": autoform.NewShortTextField().
			SetDisplayName("Website").
			SetDescription("Website URL associated with the contact").
			SetRequired(false).
			Build(),
		"address": autoform.NewShortTextField().
			SetDisplayName("Address").
			SetDescription("Street address of the contact").
			SetRequired(false).
			Build(),
		"city": autoform.NewShortTextField().
			SetDisplayName("City").
			SetDescription("City of the contact").
			SetRequired(false).
			Build(),
		"state": autoform.NewShortTextField().
			SetDisplayName("State/Region").
			SetDescription("State or region of the contact").
			SetRequired(false).
			Build(),
		"zipcode": autoform.NewShortTextField().
			SetDisplayName("Zip/Postal Code").
			SetDescription("Zip or postal code of the contact").
			SetRequired(false).
			Build(),
		"country": autoform.NewShortTextField().
			SetDisplayName("Country").
			SetDescription("Country of the contact").
			SetRequired(false).
			Build(),
	}
}

func (a *CreateContactAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createContactActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	contact := shared.ContactRequest{
		Properties: map[string]interface{}{
			"email": input.Email,
		},
	}

	if input.FirstName != "" {
		contact.Properties["firstname"] = input.FirstName
	}

	if input.LastName != "" {
		contact.Properties["lastname"] = input.LastName
	}
	if input.Phone != "" {
		contact.Properties["phone"] = input.Phone
	}
	if input.Company != "" {
		contact.Properties["company"] = input.Company
	}
	if input.JobTitle != "" {
		contact.Properties["jobtitle"] = input.JobTitle
	}
	if input.Website != "" {
		contact.Properties["website"] = input.Website
	}
	if input.Address != "" {
		contact.Properties["address"] = input.Address
	}
	if input.City != "" {
		contact.Properties["city"] = input.City
	}
	if input.State != "" {
		contact.Properties["state"] = input.State
	}
	if input.ZipCode != "" {
		contact.Properties["zip"] = input.ZipCode
	}
	if input.Country != "" {
		contact.Properties["country"] = input.Country
	}

	newContact, err := json.Marshal(contact)
	if err != nil {
		return nil, err
	}

	reqURL := "/crm/v3/objects/contacts"

	resp, err := shared.HubspotClient(reqURL, ctx.Auth.AccessToken, http.MethodPost, newContact)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *CreateContactAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateContactAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateContactAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateContactAction() sdk.Action {
	return &CreateContactAction{}
}
