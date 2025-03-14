package actions

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/wakflo/extensions/internal/integrations/keapcrm/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
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

func (a *ListContactsAction) Name() string {
	return "List Contacts"
}

func (a *ListContactsAction) Description() string {
	return "Retrieve a list of contacts from Keap with optional filtering and sorting"
}

func (a *ListContactsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListContactsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listContactsDocs,
	}
}

func (a *ListContactsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"limit": autoform.NewNumberField().
			SetDisplayName("Limit").
			SetDescription("Maximum number of contacts to return (default: 50, max: 1000)").
			Build(),
		"offset": autoform.NewNumberField().
			SetDisplayName("Offset").
			SetDescription("Number of contacts to skip (for pagination)").
			SetDefaultValue(0).
			Build(),
		"email": autoform.NewShortTextField().
			SetDisplayName("Email").
			SetDescription("Filter contacts by email address").
			SetRequired(false).Build(),
		"first_name": autoform.NewShortTextField().
			SetDisplayName("First Name").
			SetDescription("Filter contacts by first name").
			SetRequired(false).Build(),
		"last_name": autoform.NewShortTextField().
			SetDisplayName("Last Name").
			SetDescription("Filter contacts by last name").
			SetRequired(false).Build(),
		"order": autoform.NewSelectField().
			SetDisplayName("Order By").
			SetDescription("Field to order contacts by").
			SetOptions([]*sdkcore.AutoFormSchema{
				{Const: "email", Title: "Email"},
				{Const: "given_name", Title: "First Name"},
				{Const: "family_name", Title: "Last Name"},
				{Const: "last_updated", Title: "Last Updated"},
			}).
			Build(),
		"order_direction": autoform.NewSelectField().
			SetDisplayName("Order Direction").
			SetDescription("Direction of ordering").
			SetOptions([]*sdkcore.AutoFormSchema{
				{Const: "ascending", Title: "Ascending"},
				{Const: "descending", Title: "Descending"},
			}).
			SetDefaultValue("ascending").
			Build(),
	}
}

func (a *ListContactsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listContactsActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

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

	endpoint := fmt.Sprintf("/contacts?%s", queryParams.Encode())

	contactsList, err := shared.MakeKeapRequest(ctx.Auth.AccessToken, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return contactsList, nil
}

func (a *ListContactsAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListContactsAction) SampleData() sdkcore.JSON {
	return map[string]any{
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
	}
}

func (a *ListContactsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListContactsAction() sdk.Action {
	return &ListContactsAction{}
}

func (a *ListContactsAction) Icon() *string {
	icon := "mdi:account-multiple"
	return &icon
}
