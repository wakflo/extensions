package actions

import (
	"encoding/json"
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/activecampaign/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type updateContactActionProps struct {
	ContactID    string            `json:"contact-id"`
	Email        string            `json:"email"`
	FirstName    string            `json:"first-name"`
	LastName     string            `json:"last-name"`
	Phone        string            `json:"phone"`
	ListIDs      string            `json:"list-ids"`
	TagIDs       string            `json:"tag-ids"`
	CustomFields map[string]string `json:"custom-fields"`
}

type UpdateContactAction struct{}

// Metadata returns metadata about the action
func (a *UpdateContactAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "update_contact",
		DisplayName:   "Update Contact",
		Description:   "Update an existing contact in your ActiveCampaign account.",
		Type:          core.ActionTypeAction,
		Documentation: updateContactDocs,
		Icon:          "mdi:account-edit-outline",
		SampleOutput: map[string]any{
			"id":        "123",
			"email":     "updated@example.com",
			"firstName": "Updated",
			"lastName":  "User",
			"phone":     "+1234567890",
			"cdate":     "2023-01-15T15:30:00-05:00",
			"udate":     "2023-03-15T16:45:00-05:00",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *UpdateContactAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("update_contact", "Update Contact")

	form.TextField("contact-id", "contact-id").
		Placeholder("Contact ID").
		HelpText("The ID of the contact you want to update").
		Required(true)

	form.TextField("email", "email").
		Placeholder("Email").
		HelpText("Email address of the contact").
		Required(false)

	form.TextField("first-name", "first-name").
		Placeholder("First Name").
		HelpText("First name of the contact").
		Required(false)

	form.TextField("last-name", "last-name").
		Placeholder("Last Name").
		HelpText("Last name of the contact").
		Required(false)

	form.TextField("phone", "phone").
		Placeholder("Phone").
		HelpText("Phone number of the contact").
		Required(false)

	// Custom fields using the group field approach
	customFieldGroup := smartform.NewGroupFieldBuilder("field", "Custom Field")

	customFieldGroup.TextField("field", "Field ID").
		Required(true)

	customFieldGroup.TextField("value", "Value").
		Required(true)

	form.ArrayField("custom-fields", "Custom Fields").
		ItemTemplate(
			customFieldGroup.
				HelpText("Key-value pairs of custom field IDs and their values").
				Build(),
		).
		Required(false)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *UpdateContactAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *UpdateContactAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateContactActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if input.ContactID == "" {
		return nil, errors.New("contact ID is required")
	}

	contactData := map[string]interface{}{
		"contact": map[string]interface{}{},
	}
	contact := contactData["contact"].(map[string]interface{})

	if input.Email != "" {
		contact["email"] = input.Email
	}

	if input.FirstName != "" {
		contact["firstName"] = input.FirstName
	}

	if input.LastName != "" {
		contact["lastName"] = input.LastName
	}

	if input.Phone != "" {
		contact["phone"] = input.Phone
	}

	if len(input.CustomFields) > 0 {
		fieldValues := []map[string]interface{}{}

		for field, value := range input.CustomFields {
			fieldValues = append(fieldValues, map[string]interface{}{
				"field": field,
				"value": value,
			})
		}

		if len(fieldValues) > 0 {
			contact["fieldValues"] = fieldValues
		}
	}

	if len(contact) == 0 {
		return nil, errors.New("at least one field must be provided for update")
	}

	payload, err := json.Marshal(contactData)
	if err != nil {
		return nil, err
	}

	endpoint := "contacts/" + input.ContactID
	response, err := shared.PutActiveCampaignClient(
		authCtx.Extra["api_url"],
		authCtx.Extra["api_key"],
		endpoint,
		payload,
	)
	if err != nil {
		return nil, err
	}

	responseMap, ok := response.(map[string]interface{})
	if !ok {
		return nil, errors.New("unexpected response format from API")
	}

	updatedContact, ok := responseMap["contact"]
	if !ok {
		return nil, errors.New("invalid response format: contact field not found")
	}

	return updatedContact, nil
}

func NewUpdateContactAction() sdk.Action {
	return &UpdateContactAction{}
}
