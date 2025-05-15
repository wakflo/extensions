package actions

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/zohocrm/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type listRecordsActionProps struct {
	Module    string `json:"module"`
	Page      int    `json:"page"`
	PerPage   int    `json:"perPage"`
	Fields    string `json:"fields"`
	SortOrder string `json:"sortOrder"`
}

type ListRecordsAction struct{}

func (a *ListRecordsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		DisplayName:   "List Records",
		Description:   "Retrieves a list of records from a specified Zoho CRM module with pagination and sorting options",
		Type:          sdkcore.ActionTypeAction,
		Documentation: listRecordsDocs,
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
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *ListRecordsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_records", "List Records")

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

	form.NumberField("page", "Page").
		Placeholder("Page number").
		Required(false).
		DefaultValue(1).
		HelpText("Page number")

	form.NumberField("perPage", "Records Per Page").
		Placeholder("Records per page").
		Required(false).
		HelpText("Number of records to retrieve per page (max 200)")

	form.TextField("fields", "Sort Field").
		Placeholder("Sort field").
		Required(false).
		HelpText("Field to sort results by (e.g., Created_Time, Last_Name)")

	form.SelectField("sortOrder", "Sort Order").
		Placeholder("Sort order").
		Required(false).
		HelpText("Order of sorted results").
		AddOption("asc", "Ascending").
		AddOption("desc", "Descending")

	schema := form.Build()

	return schema

}

func (a *ListRecordsAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listRecordsActionProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

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

	result, err := shared.GetZohoCRMClient(token, http.MethodGet, endpoint, nil)
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

func (a *ListRecordsAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewListRecordsAction() sdk.Action {
	return &ListRecordsAction{}
}
