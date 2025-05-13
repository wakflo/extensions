package actions

import (
	"errors"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/freshworkscrm/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createNewContactActionProps struct {
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	MobileNumber *string `json:"mobile_number"`
	Email        *string `json:"email"`
	JobTitle     *string `json:"job_title,omitempty"`
	Company      *string `json:"company,omitempty"`
	Address      *string `json:"address,omitempty"`
	City         *string `json:"city,omitempty"`
	State        *string `json:"state,omitempty"`
	ZipCode      *string `json:"zip_code,omitempty"`
	Country      *string `json:"country,omitempty"`
}

type CreateNewContactAction struct{}

// Metadata returns metadata about the action
func (a *CreateNewContactAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_contact",
		DisplayName:   "Create a Contact",
		Description:   "Create a new contact in Freshworks CRM with specified information.",
		Type:          core.ActionTypeAction,
		Documentation: createContactDocs,
		SampleOutput: map[string]any{
			"contact": map[string]any{
				"id":         "12345",
				"first_name": "John",
				"last_name":  "Doe",
				"email":      "john.doe@example.com",
				"created_at": "2023-01-01T12:00:00Z",
				"updated_at": "2023-01-01T12:00:00Z",
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CreateNewContactAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_contact", "Create a Contact")

	// Add first_name field
	form.TextField("first_name", "First Name").
		Placeholder("Enter first name").
		Required(true).
		HelpText("Contact's first name")

	// Add last_name field
	form.TextField("last_name", "Last Name").
		Placeholder("Enter last name").
		Required(true).
		HelpText("Contact's last name")

	// Add email field
	form.TextField("email", "Email").
		Placeholder("Enter email").
		Required(true).
		HelpText("Contact's email")

	// Add mobile_number field
	form.TextField("mobile_number", "Mobile Number").
		Placeholder("Enter mobile number").
		Required(false).
		HelpText("Contact's mobile number")

	// Add job_title field
	form.TextField("job_title", "Job Title").
		Placeholder("Enter job title").
		Required(false).
		HelpText("Job title of the contact")

	// Add company field
	form.TextField("company", "Company").
		Placeholder("Enter company").
		Required(false).
		HelpText("Company name of the contact")

	// Add address field
	form.TextField("address", "Address").
		Placeholder("Enter address").
		Required(false).
		HelpText("Street address of the contact")

	// Add city field
	form.TextField("city", "City").
		Placeholder("Enter city").
		Required(false).
		HelpText("City of the contact")

	// Add state field
	form.TextField("state", "State").
		Placeholder("Enter state").
		Required(false).
		HelpText("State/Province of the contact")

	// Add zip_code field
	form.TextField("zip_code", "Zip Code").
		Placeholder("Enter zip code").
		Required(false).
		HelpText("Postal/Zip code of the contact")

	// Add country field
	form.TextField("country", "Country").
		Placeholder("Enter country").
		Required(false).
		HelpText("Country of the contact")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CreateNewContactAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CreateNewContactAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["api-key"] == "" || authCtx.Extra["domain"] == "" {
		return nil, errors.New("missing freshworks auth parameters")
	}

	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[createNewContactActionProps](ctx)
	if err != nil {
		return nil, err
	}

	domain := authCtx.Extra["domain"]
	freshworksDomain := "https://" + domain + ".myfreshworks.com"

	contactData := map[string]interface{}{
		"contact": map[string]interface{}{
			"first_name":    input.FirstName,
			"last_name":     input.LastName,
			"email":         input.Email,
			"mobile_number": input.MobileNumber,
		},
	}

	contactMap := contactData["contact"].(map[string]interface{})
	if input.MobileNumber != nil {
		contactMap["phone"] = *input.MobileNumber
	}
	if input.JobTitle != nil {
		contactMap["job_title"] = *input.JobTitle
	}
	if input.Company != nil {
		contactMap["company"] = *input.Company
	}
	if input.Address != nil {
		contactMap["address"] = *input.Address
	}
	if input.City != nil {
		contactMap["city"] = *input.City
	}
	if input.State != nil {
		contactMap["state"] = *input.State
	}
	if input.ZipCode != nil {
		contactMap["zipcode"] = *input.ZipCode
	}
	if input.Country != nil {
		contactMap["country"] = *input.Country
	}

	response, err := shared.CreateContact(freshworksDomain, authCtx.Extra["api-key"], contactData)
	if err != nil {
		return nil, fmt.Errorf("error creating contact:  %v", err)
	}

	return response, nil
}

func NewCreateNewContactAction() sdk.Action {
	return &CreateNewContactAction{}
}
