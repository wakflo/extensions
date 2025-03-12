package actions

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/wakflo/extensions/internal/integrations/zohocrm/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listRecordsActionProps struct {
	Module    string `json:"module"`
	Page      int    `json:"page"`
	PerPage   int    `json:"perPage"`
	Fields    string `json:"fields"`
	SortOrder string `json:"sortOrder"`
}

type ListRecordsAction struct{}

func (a *ListRecordsAction) Name() string {
	return "List Records"
}

func (a *ListRecordsAction) Description() string {
	return "Retrieves a list of records from a specified Zoho CRM module with pagination and sorting options"
}

func (a *ListRecordsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListRecordsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listRecordsDocs,
	}
}

func (a *ListRecordsAction) Icon() *string {
	return nil
}

func (a *ListRecordsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"module": autoform.NewDynamicField(sdkcore.String).
			SetDisplayName("Module").
			SetDescription("The Zoho CRM module to list records from (e.g., Leads, Contacts, Accounts)").
			SetDynamicOptions(shared.GetModulesFunction()).
			SetRequired(true).
			Build(),
		"page": autoform.NewNumberField().
			SetDisplayName("Page").
			SetDescription("Page number for pagination (starts from 1)").
			SetDefaultValue(1).
			Build(),
		"perPage": autoform.NewNumberField().
			SetDisplayName("Records Per Page").
			SetDescription("Number of records to retrieve per page (max 200)").
			Build(),
		"fields": autoform.NewShortTextField().
			SetDisplayName("Sort Field").
			SetDescription("Field to sort results by (e.g., Created_Time, Last_Name)").
			Build(),
		"sortOrder": autoform.NewSelectField().
			SetDisplayName("Sort Order").
			SetDescription("Order of sorted results").
			SetOptions([]*sdkcore.AutoFormSchema{
				{
					Title: "Ascending",
					Const: "asc",
				},
				{
					Title: "Descending",
					Const: "desc",
				},
			}).
			SetDefaultValue("asc").
			Build(),
	}
}

func (a *ListRecordsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listRecordsActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	// Build query parameters
	queryParams := url.Values{}

	// Set required fields parameter
	if input.Fields != "" {
		queryParams.Set("fields", input.Fields)
	} else {
		// Fallback to basic fields if empty
		queryParams.Set("fields", "id,Created_Time,Modified_Time")
	}

	// Set pagination parameters
	if input.Page > 0 {
		page := input.Page - 1 // Zoho uses 0-indexed pages
		queryParams.Set("page", strconv.Itoa(page))
	}

	if input.PerPage > 0 {
		queryParams.Set("per_page", strconv.Itoa(input.PerPage))
	}

	endpoint := input.Module
	if len(queryParams) > 0 {
		endpoint = fmt.Sprintf("%s?%s", endpoint, queryParams.Encode())
	}

	result, err := shared.GetZohoCRMClient(ctx.Auth.AccessToken, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error calling Zoho CRM API: %v", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok {
		// Check if data is nil, which could mean no records found
		if result["data"] == nil {
			return map[string]interface{}{
				"records": []interface{}{},
				"count":   0,
				"module":  input.Module,
			}, nil
		}
		return nil, errors.New("invalid response format: data field is missing or not an array")
	}

	info, infoOk := result["info"].(map[string]interface{})

	response := map[string]interface{}{
		"records": data,
		"count":   len(data),
		"module":  input.Module,
	}

	// Add pagination info if available
	if infoOk {
		response["info"] = info
	}

	return response, nil
}

func (a *ListRecordsAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListRecordsAction) SampleData() sdkcore.JSON {
	return map[string]interface{}{
		"records": []interface{}{
			map[string]interface{}{
				"id":          "3477061000000419001",
				"Last_Name":   "Smith",
				"Email":       "john.smith@example.com",
				"Company":     "ACME Corp",
				"Lead_Status": "Qualified",
			},
			map[string]interface{}{
				"id":          "3477061000000419002",
				"Last_Name":   "Johnson",
				"Email":       "sarah.johnson@example.com",
				"Company":     "XYZ Inc",
				"Lead_Status": "New",
			},
		},
		"info": map[string]interface{}{
			"per_page":     "100",
			"count":        "2",
			"page":         "1",
			"more_records": false,
		},
	}
}

func (a *ListRecordsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListRecordsAction() sdk.Action {
	return &ListRecordsAction{}
}
