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

type updateContactActionProps struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	MobileNumber  string `json:"mobile_number"`
	Email         string `json:"email"`
	ContactViewID string `json:"contact_view_id"`
	ContactID     string `json:"contact_id"`
}

type UpdateContactAction struct{}

// Metadata returns metadata about the action
func (a *UpdateContactAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "update_contact",
		DisplayName:   "Update a Contact",
		Description:   "Update a contact in Freshworks CRM with specified information.",
		Type:          core.ActionTypeAction,
		Documentation: updateContactDocs,
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
func (a *UpdateContactAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("update_contact", "Update a Contact")

	// Add contact_view_id field
	form.TextField("contact_view_id", "Contact View ID").
		Placeholder("Enter contact view ID").
		Required(true).
		HelpText("The ID of the contact view")

	// Add contact_id field
	form.TextField("contact_id", "Contact ID").
		Placeholder("Enter contact ID").
		Required(true).
		HelpText("The ID of the contact to update")

	// Add first_name field
	form.TextField("first_name", "First Name").
		Placeholder("Enter first name").
		Required(false).
		HelpText("Contact's first name")

	// Add last_name field
	form.TextField("last_name", "Last Name").
		Placeholder("Enter last name").
		Required(false).
		HelpText("Contact's last name")

	// Add email field
	form.TextField("email", "Email").
		Placeholder("Enter email").
		Required(false).
		HelpText("Contact's email")

	// Add mobile_number field
	form.TextField("mobile_number", "Mobile Number").
		Placeholder("Enter mobile number").
		Required(false).
		HelpText("Contact's mobile number")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *UpdateContactAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *UpdateContactAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["api-key"] == "" || authCtx.Extra["domain"] == "" {
		return nil, errors.New("missing freshworks auth parameters")
	}

	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[updateContactActionProps](ctx)
	if err != nil {
		return nil, err
	}

	domain := authCtx.Extra["domain"]
	freshworksDomain := "https://" + domain + ".myfreshworks.com"

	contact := make(map[string]interface{})

	shared.UpdateField(contact, "first_name", input.FirstName)
	shared.UpdateField(contact, "last_name", input.LastName)
	shared.UpdateField(contact, "email", input.Email)
	shared.UpdateField(contact, "mobile_number", input.MobileNumber)

	contactData := map[string]interface{}{
		"contact": contact,
	}

	response, err := shared.UpdateContact(freshworksDomain, authCtx.Extra["api-key"], input.ContactID, contactData)
	if err != nil {
		return nil, fmt.Errorf("error updating contact:  %v", err)
	}

	return response, nil
}

func NewUpdateContactAction() sdk.Action {
	return &UpdateContactAction{}
}
