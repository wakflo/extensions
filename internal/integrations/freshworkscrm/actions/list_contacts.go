package actions

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/wakflo/extensions/internal/integrations/freshworkscrm/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listContactsActionProps struct {
	Page     int    `json:"page"`
	PerPage  int    `json:"per_page"`
	SortBy   string `json:"sort_by"`
	FilterBy string `json:"filter_by"`
}

type ListContactsAction struct{}

func (a *ListContactsAction) Name() string {
	return "List Contacts"
}

func (a *ListContactsAction) Description() string {
	return "Retrieve a list of contacts from Freshworks CRM with pagination and filtering options."
}

func (a *ListContactsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListContactsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listContactsDocs,
	}
}

func (a *ListContactsAction) Icon() *string {
	return nil
}

func (a *ListContactsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"page": autoform.NewNumberField().
			SetDisplayName("Page").
			SetDescription("Page number for pagination").
			SetRequired(false).
			Build(),
		"per_page": autoform.NewNumberField().
			SetDisplayName("Per Page").
			SetDescription("Number of contacts per page").
			SetRequired(false).
			Build(),
		"sort_by": autoform.NewShortTextField().
			SetDisplayName("Sort By").
			SetDescription("Field to sort contacts by (e.g., first_name, last_name, created_at)").
			SetRequired(false).
			Build(),
		"filter_by": autoform.NewShortTextField().
			SetDisplayName("Filter By").
			SetDescription("Filter contacts by attributes (e.g., email=test@example.com)").
			SetRequired(false).
			Build(),
	}
}

func (a *ListContactsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" || ctx.Auth.Extra["domain"] == "" {
		return nil, errors.New("missing freshworks auth parameters")
	}

	input := sdk.InputToType[listContactsActionProps](ctx.BaseContext)

	domain := ctx.Auth.Extra["domain"]
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

	response, err := shared.ListContacts(freshworksDomain, ctx.Auth.Extra["api-key"], queryParams)
	if err != nil {
		return nil, fmt.Errorf("error listing contacts: %v", err)
	}

	return response, nil
}

func (a *ListContactsAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListContactsAction) SampleData() sdkcore.JSON {
	return map[string]any{
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
	}
}

func (a *ListContactsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListContactsAction() sdk.Action {
	return &ListContactsAction{}
}
