package actions

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/zohocrm/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
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

func (a *SearchRecordsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "search_records",
		DisplayName:   "Search Records",
		Description:   "Searches for records in a Zoho CRM module based on specified criteria",
		Type:          sdkcore.ActionTypeAction,
		Documentation: searchRecordsDocs,
		SampleOutput: map[string]interface{}{
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
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *SearchRecordsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("search_records", "Search Records")

	form.SelectField("module", "Module").
		Placeholder("Select a module").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(shared.GetModulesFunction())).
				WithSearchSupport().
				End().GetDynamicSource(),
		)

	form.TextField("criteria", "Search Criteria").
		Placeholder("Last_Name:equals:Smith").
		Required(true).
		HelpText("Search criteria in the format: field:operator:value. Example: Last_Name:equals:Smith or Email:contains:example.com")

	form.NumberField("page", "Page").
		Placeholder("Page number").
		Required(false).
		DefaultValue(1).
		HelpText("Page number")

	form.NumberField("perPage", "Records Per Page").
		Placeholder("Records per page").
		Required(false).
		HelpText("Number of records to retrieve per page (max 200)")

	form.SelectField("sortField", "Sort Field").
		Placeholder("Sort field").
		Required(false).
		HelpText("Field to sort results by (e.g., Created_Time, Last_Name)").
		AddOption("", "None").
		AddOption("id", "ID").
		AddOption("Created_Time", "Created Time").
		AddOption("Modified_Time", "Modified Time").
		DefaultValue("")

	form.SelectField("sortOrder", "Sort Order").
		Placeholder("Sort order").
		Required(false).
		HelpText("Order of sorted results").
		AddOption("asc", "Ascending").
		AddOption("desc", "Descending")

	schema := form.Build()

	return schema
}

func (a *SearchRecordsAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[searchRecordsActionProps](ctx)
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

	result, err := shared.GetZohoCRMClient(ctx.Auth().AccessToken, http.MethodGet, endpoint, nil)
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

func (a *SearchRecordsAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewSearchRecordsAction() sdk.Action {
	return &SearchRecordsAction{}
}
