package actions

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/freshworkscrm/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type listContactsActionProps struct {
	Page     int    `json:"page"`
	PerPage  int    `json:"per_page"`
	SortBy   string `json:"sort_by"`
	FilterBy string `json:"filter_by"`
}

type ListContactsAction struct{}

// Metadata returns metadata about the action
func (a *ListContactsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_contacts",
		DisplayName:   "List Contacts",
		Description:   "Retrieve a list of contacts from Freshworks CRM with pagination and filtering options.",
		Type:          core.ActionTypeAction,
		Documentation: listContactsDocs,
		SampleOutput: map[string]any{
			"contacts": []map[string]any{
				{
					"id":            "12345",
					"first_name":    "John",
					"last_name":     "Doe",
					"email":         "john.doe@example.com",
					"mobile_number": "+1234567890",
					"created_at":    "2023-01-01T12:00:00Z",
					"updated_at":    "2023-01-01T12:00:00Z",
				},
				{
					"id":            "12346",
					"first_name":    "Jane",
					"last_name":     "Smith",
					"email":         "jane.smith@example.com",
					"mobile_number": "+0987654321",
					"created_at":    "2023-01-02T12:00:00Z",
					"updated_at":    "2023-01-02T12:00:00Z",
				},
			},
			"meta": map[string]any{
				"total_pages": "10",
				"total_count": "245",
				"per_page":    "25",
				"page":        "1",
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ListContactsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_contacts", "List Contacts")

	// Add page field
	form.NumberField("page", "Page").
		Placeholder("Enter page number").
		Required(false).
		HelpText("Page number for pagination")

	// Add per_page field
	form.NumberField("per_page", "Per Page").
		Placeholder("Enter items per page").
		Required(false).
		HelpText("Number of contacts per page")

	// Add sort_by field
	form.TextField("sort_by", "Sort By").
		Placeholder("Enter sort field").
		Required(false).
		HelpText("Field to sort contacts by (e.g., first_name, last_name, created_at)")

	// Add filter_by field
	form.TextField("filter_by", "Filter By").
		Placeholder("Enter filter criteria").
		Required(false).
		HelpText("Filter contacts by attributes (e.g., email=test@example.com)")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ListContactsAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ListContactsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["api-key"] == "" || authCtx.Extra["domain"] == "" {
		return nil, errors.New("missing freshworks auth parameters")
	}

	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[listContactsActionProps](ctx)
	if err != nil {
		return nil, err
	}

	domain := authCtx.Extra["domain"]
	freshworksDomain := "https://" + domain + ".myfreshworks.com"

	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PerPage <= 0 {
		input.PerPage = 25
	}

	queryParams := map[string]string{
		"page":     strconv.Itoa(input.Page),
		"per_page": strconv.Itoa(input.PerPage),
	}

	if input.FilterBy != "" {
		queryParams["filter"] = input.FilterBy
	}

	response, err := shared.ListContacts(freshworksDomain, authCtx.Extra["api-key"], queryParams)
	if err != nil {
		return nil, fmt.Errorf("error listing contacts: %v", err)
	}

	return response, nil
}

func NewListContactsAction() sdk.Action {
	return &ListContactsAction{}
}
