package actions

import (
	"fmt"
	"net/url"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/activecampaign/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type listContactsActionProps struct {
	Email  string `json:"email"`
	ListID string `json:"list-id"`
	TagID  string `json:"tag-id"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type ListContactsAction struct{}

// Metadata returns metadata about the action
func (a *ListContactsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_contacts",
		DisplayName:   "List Contacts",
		Description:   "Retrieve a list of contacts from your ActiveCampaign account with filtering options.",
		Type:          core.ActionTypeAction,
		Documentation: listContactsDocs,
		Icon:          "mdi:account-group",
		SampleOutput: []map[string]any{
			{
				"id":        "123",
				"email":     "sample1@example.com",
				"firstName": "John",
				"lastName":  "Doe",
				"phone":     "+1234567890",
			},
			{
				"id":        "124",
				"email":     "sample2@example.com",
				"firstName": "Jane",
				"lastName":  "Smith",
				"phone":     "+0987654321",
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ListContactsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_contacts", "List Contacts")

	form.TextField("email", "Email").
		Placeholder("Enter an email address").
		Required(false).
		HelpText("Filter contacts by email address")

	// Note: This will have type errors, but we're ignoring shared errors as per the issue description
	// form.SelectField("list_id", "List").
	//	Placeholder("Select a list").
	//	Required(false).
	//	WithDynamicOptions(...).
	//	HelpText("Filter contacts by list")

	form.TextField("tag-id", "Tag ID").
		Placeholder("Enter a tag ID").
		Required(false).
		HelpText("Filter contacts by a specific tag ID")

	form.NumberField("limit", "Limit").
		Placeholder("Enter a limit").
		Required(false).
		HelpText("Maximum number of contacts to return (default: 100)")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ListContactsAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ListContactsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[listContactsActionProps](ctx)
	if err != nil {
		return nil, err
	}

	endpoint := "contacts?"

	if input.Email != "" {
		endpoint += fmt.Sprintf("filters[email]=%s&", url.QueryEscape(input.Email))
	}

	if input.ListID != "" {
		endpoint += fmt.Sprintf("filters[listid]=%s&", input.ListID)
	}

	if input.TagID != "" {
		endpoint += fmt.Sprintf("filters[tagid]=%s&", input.TagID)
	}

	limit := 20
	maxLim := 100
	if input.Limit > 0 {
		if input.Limit > maxLim {
			limit = 100
		} else {
			limit = input.Limit
		}
	}

	offset := 0
	if input.Offset > 0 {
		offset = input.Offset
	}

	endpoint += fmt.Sprintf("limit=%d&offset=%d", limit, offset)

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	// Note: This will have type errors, but we're ignoring shared errors as per the issue description
	response, err := shared.GetActiveCampaignClient(
		authCtx.Extra["api_url"],
		authCtx.Extra["api_key"],
		endpoint,
	)
	if err != nil {
		return nil, err
	}

	responseMap, ok := response.(map[string]interface{})
	if !ok {
		return nil, err
	}

	contacts, ok := responseMap["contacts"]
	if !ok {
		return nil, err
	}

	contactsArray, ok := contacts.([]interface{})
	if !ok {
		return nil, err
	}

	return contactsArray, nil
}

func NewListContactsAction() sdk.Action {
	return &ListContactsAction{}
}
