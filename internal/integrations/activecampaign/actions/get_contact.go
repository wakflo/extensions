package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/activecampaign/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getContactActionProps struct {
	ContactID string `json:"contact-id"`
}

type GetContactAction struct{}

// Metadata returns metadata about the action
func (a *GetContactAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_contact",
		DisplayName:   "Get Contact",
		Description:   "Retrieve a specific contact by ID from your ActiveCampaign account.",
		Type:          core.ActionTypeAction,
		Documentation: getContactDocs,
		Icon:          "mdi:account-search",
		SampleOutput: map[string]any{
			"id":        "123",
			"email":     "sample@example.com",
			"firstName": "John",
			"lastName":  "Doe",
			"phone":     "+1234567890",
			"cdate":     "2023-01-15T15:30:00-05:00",
			"udate":     "2023-02-20T10:15:00-05:00",
			"links": map[string]string{
				"lists": "https://api.example.com/contacts/123/lists",
				"deals": "https://api.example.com/contacts/123/deals",
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GetContactAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_contact", "Get Contact")

	form.TextField("contact-id", "contact-id").
		Placeholder("Contact ID").
		HelpText("The ID of the contact you want to retrieve").
		Required(true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetContactAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GetContactAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getContactActionProps](ctx)
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

	endpoint := "contacts/" + input.ContactID

	response, err := shared.GetActiveCampaignClient(
		authCtx.Extra["api_url"],
		authCtx.Extra["api_key"],
		endpoint,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewGetContactAction() sdk.Action {
	return &GetContactAction{}
}
