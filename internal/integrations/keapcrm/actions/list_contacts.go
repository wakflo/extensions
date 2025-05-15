package actions

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/keapcrm/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type listContactsActionProps struct {
	Limit          int    `json:"limit,omitempty"`
	Offset         int    `json:"offset,omitempty"`
	Email          string `json:"email,omitempty"`
	FirstName      string `json:"first_name,omitempty"`
	LastName       string `json:"last_name,omitempty"`
	Order          string `json:"order,omitempty"`
	OrderDirection string `json:"order_direction,omitempty"`
}

type ListContactsAction struct{}

func (a *ListContactsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		DisplayName:   "List Contacts",
		Description:   "Retrieve a list of contacts from Keap with optional filtering and sorting",
		Type:          sdkcore.ActionTypeAction,
		Documentation: listContactsDocs,
		SampleOutput: map[string]any{
			"contacts": []map[string]any{
				{
					"id":          "12345",
					"given_name":  "John",
					"family_name": "Doe",
					"email_addresses": []map[string]string{
						{
							"email": "john.doe@example.com",
							"type":  "PRIMARY",
						},
					},
					"last_updated": "2023-01-15T10:30:00Z",
				},
				{
					"id":          "67890",
					"given_name":  "Jane",
					"family_name": "Smith",
					"email_addresses": []map[string]string{
						{
							"email": "jane.smith@example.com",
							"type":  "PRIMARY",
						},
					},
					"last_updated": "2023-02-20T14:45:00Z",
				},
			},
			"total":  "2",
			"limit":  "50",
			"offset": "0",
		},

		Settings: sdkcore.ActionSettings{},
	}
}

func (a *ListContactsAction) Properties() *smartform.FormSchema {

	form := smartform.NewForm("list_contacts", "List Contacts")

	form.NumberField("limit", "Limit").
		Placeholder("Enter a value for Limit.").
		Required(false).
		HelpText("Maximum number of contacts to return (default: 50, max: 1000).")

	form.NumberField("offset", "Offset").
		Placeholder("Enter a value for Offset.").
		Required(false).
		HelpText("Number of contacts to skip (for pagination).")

	form.TextField("email", "Email").
		Placeholder("Enter a value for Email.").
		Required(false).
		HelpText("Filter contacts by email address.")

	form.TextField("first_name", "First Name").
		Placeholder("Enter a value for First Name.").
		Required(false).
		HelpText("Filter contacts by first name.")

	form.TextField("last_name", "Last Name").
		Placeholder("Enter a value for Last Name.").
		Required(false).
		HelpText("Filter contacts by last name.")

	form.SelectField("order", "Order By").
		Placeholder("Select a value for Order By.").
		Required(false).
		HelpText("Field to order contacts by.").
		AddOptions([]*smartform.Option{
			{Value: "email", Label: "Email"},
			{Value: "given", Label: "First Name"},
			{Value: "family", Label: "Last Name"},
			{Value: "last_updated", Label: "Last Updated"},
		}...)

	form.SelectField("order_direction", "Order Direction").
		Placeholder("Select a value for Order Direction.").
		Required(false).
		HelpText("Direction of ordering.").
		AddOptions([]*smartform.Option{
			{Value: "ascending", Label: "Ascending"},
			{Value: "descending", Label: "Descending"}}...).
		DefaultValue("ascending")

	schema := form.Build()

	return schema

}

func (a *ListContactsAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listContactsActionProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	// Prepare query parameters
	queryParams := url.Values{}

	// Add pagination parameters
	if input.Limit > 0 {
		queryParams.Add("limit", strconv.Itoa(input.Limit))
	} else {
		queryParams.Add("limit", "50") // Default limit
	}

	if input.Offset > 0 {
		queryParams.Add("offset", strconv.Itoa(input.Offset))
	}

	if input.Email != "" {
		queryParams.Add("email", input.Email)
	}
	if input.FirstName != "" {
		queryParams.Add("given_name", input.FirstName)
	}
	if input.LastName != "" {
		queryParams.Add("family_name", input.LastName)
	}

	if input.Order != "" {
		queryParams.Add("order", input.Order)
	}
	if input.OrderDirection != "" {
		queryParams.Add("order_direction", input.OrderDirection)
	}

	endpoint := "/contacts?" + queryParams.Encode()

	contactsList, err := shared.MakeKeapRequest(token, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return contactsList, nil
}

func (a *ListContactsAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewListContactsAction() sdk.Action {
	return &ListContactsAction{}
}
