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

type createContactActionProps struct {
	Email     string `json:"email"`
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
	Phone     string `json:"phone"`
}

type CreateContactAction struct{}

// Metadata returns metadata about the action
func (a *CreateContactAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_contact",
		DisplayName:   "Create Contact",
		Description:   "Create a new contact in your ActiveCampaign account with customizable fields.",
		Type:          core.ActionTypeAction,
		Documentation: createContactDocs,
		Icon:          "mdi:account-plus",
		SampleOutput: map[string]any{
			"id":        "123",
			"email":     "new@example.com",
			"firstName": "New",
			"lastName":  "User",
			"phone":     "+1234567890",
			"cdate":     "2023-03-15T15:30:00-05:00",
			"udate":     "2023-03-15T15:30:00-05:00",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CreateContactAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_contact", "Create Contact")

	form.TextField("email", "email").
		Placeholder("Email").
		HelpText("Email address of the contact").
		Required(true)

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

	if input.Email == "" {
		return nil, errors.New("email is required")
	}

	contactData := map[string]interface{}{
		"contact": map[string]interface{}{
			"email": input.Email,
		},
	}

	contact := contactData["contact"].(map[string]interface{})

	if input.FirstName != "" {
		contact["firstName"] = input.FirstName
	}

	if input.LastName != "" {
		contact["lastName"] = input.LastName
	}

	if input.Phone != "" {
		contact["phone"] = input.Phone
	}

	payload, err := json.Marshal(contactData)
	if err != nil {
		return nil, err
	}

	response, err := shared.PostActiveCampaignClient(
		authCtx.Extra["api_url"],
		authCtx.Extra["api_key"],
		"contacts",
		payload,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewCreateContactAction() sdk.Action {
	return &CreateContactAction{}
}
