package actions

import (
	"encoding/json"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/hubspot/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
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

// Metadata returns metadata about the action
func (a *CreateContactAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_contact",
		DisplayName:   "Create Contact",
		Description:   "Create a new contact in HubSpot with the specified information.",
		Type:          core.ActionTypeAction,
		Documentation: createContactDocs,
		SampleOutput: map[string]any{
			"id": "12345",
			"properties": map[string]any{
				"email":     "john.doe@example.com",
				"firstname": "John",
				"lastname":  "Doe",
				"phone":     "+1234567890",
				"company":   "Example Inc.",
			},
			"createdAt": "2023-05-01T12:00:00.000Z",
			"updatedAt": "2023-05-01T12:00:00.000Z",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CreateContactAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_contact", "Create Contact")

	form.TextField("email", "Email").
		Required(true).
		HelpText("Email address of the contact")

	form.TextField("firstname", "First Name").
		Required(false).
		HelpText("First name of the contact")

	form.TextField("lastname", "Last Name").
		Required(false).
		HelpText("Last name of the contact")

	form.TextField("phone", "Phone").
		Required(false).
		HelpText("Phone number of the contact")

	form.TextField("company", "Company").
		Required(false).
		HelpText("Company name the contact works for")

	form.TextField("jobtitle", "Job Title").
		Required(false).
		HelpText("Job title of the contact")

	form.TextField("website", "Website").
		Required(false).
		HelpText("Website URL associated with the contact")

	form.TextareaField("address", "Address").
		Required(false).
		HelpText("Street address of the contact")

	form.TextField("city", "City").
		Required(false).
		HelpText("City of the contact")

	form.TextField("state", "State/Region").
		Required(false).
		HelpText("State or region of the contact")

	form.TextField("zipcode", "Zip/Postal Code").
		Required(false).
		HelpText("Zip or postal code of the contact")

	form.TextField("country", "Country").
		Required(false).
		HelpText("Country of the contact")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CreateContactAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CreateContactAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createContactActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
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

	resp, err := shared.HubspotClient(reqURL, authCtx.Token.AccessToken, http.MethodPost, newContact)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewCreateContactAction() sdk.Action {
	return &CreateContactAction{}
}
