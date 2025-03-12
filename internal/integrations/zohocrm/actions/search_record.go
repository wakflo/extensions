package actions

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/wakflo/extensions/internal/integrations/zohocrm/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type searchRecordsActionProps struct {
	Module    string `json:"module"`
	Criteria  string `json:"criteria"`
	Page      int    `json:"page"`
	PerPage   int    `json:"perPage"`
	SortField string `json:"sortField"`
	SortOrder string `json:"sortOrder"`
}

type SearchRecordsAction struct{}

func (a *SearchRecordsAction) Name() string {
	return "Search Records"
}

func (a *SearchRecordsAction) Description() string {
	return "Searches for records in a Zoho CRM module based on specified criteria"
}

func (a *SearchRecordsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *SearchRecordsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &searchRecordsDocs,
	}
}

func (a *SearchRecordsAction) Icon() *string {
	return nil
}

func (a *SearchRecordsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"module": autoform.NewDynamicField(sdkcore.String).
			SetDisplayName("Module").
			SetDescription("The Zoho CRM module to search records in (e.g., Leads, Contacts, Accounts)").
			SetDynamicOptions(shared.GetModulesFunction()).
			SetRequired(true).
			SetPlaceholder("Leads").
			Build(),
		"criteria": autoform.NewShortTextField().
			SetDisplayName("Search Criteria").
			SetDescription("Search criteria in the format: field:operator:value. Example: Last_Name:equals:Smith or Email:contains:example.com").
			SetRequired(true).
			SetPlaceholder("Last_Name:equals:Smith").
			Build(),
		"page": autoform.NewNumberField().
			SetDisplayName("Page").
			SetDescription("Page number for pagination (starts from 1)").
			SetPlaceholder("1").
			Build(),
		"perPage": autoform.NewNumberField().
			SetDisplayName("Records Per Page").
			SetDescription("Number of records to retrieve per page (max 200)").
			SetPlaceholder("100").
			Build(),
		"sortField": autoform.NewSelectField().
			SetDisplayName("Sort Field").
			SetDescription("Field to sort results by (only certain fields are supported by Zoho CRM)").
			SetOptions([]*sdkcore.AutoFormSchema{
				{
					Title: "None",
					Const: "",
				},
				{
					Title: "ID",
					Const: "id",
				},
				{
					Title: "Created Time",
					Const: "Created_Time",
				},
				{
					Title: "Modified Time",
					Const: "Modified_Time",
				},
			}).
			SetDefaultValue("").
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

func (a *SearchRecordsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[searchRecordsActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	queryParams := url.Values{}

	if input.Criteria != "" {
		criteria := strings.Trim(input.Criteria, "()")

		const lenCriteria = 3

		parts := strings.Split(criteria, ":")
		if len(parts) < lenCriteria {
			return nil, errors.New("invalid criteria format: must be field:operator:value")
		}

		formattedCriteria := fmt.Sprintf("(%s)", criteria)
		queryParams.Set("criteria", formattedCriteria)
	}

	if input.Page > 0 {
		page := input.Page - 1
		queryParams.Set("page", strconv.Itoa(page))
	}

	if input.PerPage > 0 {
		perPage := input.PerPage
		queryParams.Set("per_page", strconv.Itoa(perPage))
	}

	if input.SortField != "" {
		validSortFields := map[string]bool{
			"id":            true,
			"Created_Time":  true,
			"Modified_Time": true,
		}

		if validSortFields[input.SortField] {
			sortOrder := "asc"
			if input.SortOrder == "desc" {
				sortOrder = "desc"
			}
			queryParams.Set("sort_by", input.SortField)
			queryParams.Set("sort_order", sortOrder)
		} else {
			return nil, errors.New("invalid sort field. Zoho CRM only supports sorting by 'id', 'Created_Time', or 'Modified_Time'")
		}
	}

	endpoint := input.Module + "/search"
	if len(queryParams) > 0 {
		endpoint = endpoint + "?" + queryParams.Encode()
	}

	result, err := shared.GetZohoCRMClient(ctx.Auth.AccessToken, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error calling Zoho CRM API: %v", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok {
		// Handle case where data is nil (no results)
		if emptyData, exists := result["data"]; exists && emptyData == nil {
			return map[string]interface{}{
				"records": []interface{}{},
				"info": map[string]interface{}{
					"per_page":     input.PerPage,
					"count":        0,
					"page":         input.Page,
					"more_records": false,
				},
			}, nil
		}
		return nil, errors.New("invalid response format: data field is missing or not an array")
	}

	info, infoOk := result["info"].(map[string]interface{})

	response := map[string]interface{}{
		"records": data,
	}

	if infoOk {
		response["info"] = info
	}

	return response, nil
}

func (a *SearchRecordsAction) Auth() *sdk.Auth {
	return nil
}

func (a *SearchRecordsAction) SampleData() sdkcore.JSON {
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
				"id":          "3477061000000419005",
				"Last_Name":   "Smith",
				"Email":       "robert.smith@example.com",
				"Company":     "ABC Ltd",
				"Lead_Status": "Contacted",
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

func (a *SearchRecordsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewSearchRecordsAction() sdk.Action {
	return &SearchRecordsAction{}
}
