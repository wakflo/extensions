package actions

import (
	"encoding/json"
	"errors"

	"github.com/wakflo/extensions/internal/integrations/activecampaign/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
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

func (a *UpdateContactAction) Name() string {
	return "Update Contact"
}

func (a *UpdateContactAction) Description() string {
	return "Update an existing contact in your ActiveCampaign account."
}

func (a *UpdateContactAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *UpdateContactAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &updateContactDocs,
	}
}

func (a *UpdateContactAction) Icon() *string {
	icon := "mdi:account-edit-outline"
	return &icon
}

func (a *UpdateContactAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"contact-id": autoform.NewShortTextField().
			SetDisplayName("Contact ID").
			SetDescription("The ID of the contact you want to update").
			SetRequired(true).
			Build(),
		"email": autoform.NewShortTextField().
			SetDisplayName("Email").
			SetDescription("Email address of the contact").
			SetRequired(false).
			Build(),
		"first-name": autoform.NewShortTextField().
			SetDisplayName("First Name").
			SetDescription("First name of the contact").
			SetRequired(false).
			Build(),
		"last-name": autoform.NewShortTextField().
			SetDisplayName("Last Name").
			SetDescription("Last name of the contact").
			SetRequired(false).
			Build(),
		"phone": autoform.NewShortTextField().
			SetDisplayName("Phone").
			SetDescription("Phone number of the contact").
			SetRequired(false).
			Build(),
		"custom-fields": autoform.NewObjectField().
			SetProperties(map[string]*sdkcore.AutoFormSchema{
				"field": autoform.NewShortTextField().
					SetDisplayName("Custom Field ID").
					SetDescription("ID of the custom field").
					SetRequired(true).
					Build(),
				"value": autoform.NewShortTextField().
					SetDisplayName("Custom Field ID").
					SetDescription("ID of the custom field").
					SetRequired(true).
					Build(),
			}).
			SetDisplayName("Custom Fields").
			SetDescription("Key-value pairs of custom field IDs and their values").
			SetRequired(false).
			Build(),
	}
}

func (a *UpdateContactAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateContactActionProps](ctx.BaseContext)
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
		ctx.Auth.Extra["api_url"],
		ctx.Auth.Extra["api_key"],
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

func (a *UpdateContactAction) Auth() *sdk.Auth {
	return nil
}

func (a *UpdateContactAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":        "123",
		"email":     "updated@example.com",
		"firstName": "Updated",
		"lastName":  "User",
		"phone":     "+1234567890",
		"cdate":     "2023-01-15T15:30:00-05:00",
		"udate":     "2023-03-15T16:45:00-05:00",
	}
}

func (a *UpdateContactAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewUpdateContactAction() sdk.Action {
	return &UpdateContactAction{}
}
