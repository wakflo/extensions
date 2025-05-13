package triggers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/zohocrm/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type recordCreatedTriggerProps struct {
	Module string `json:"module"`
}

type RecordCreatedTrigger struct{}

func (t *RecordCreatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "record_created",
		DisplayName:   "Record Created",
		Description:   "Triggers when a new record is created in a specified Zoho CRM module",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: recordCreatedDocs,
		SampleOutput: map[string]interface{}{
			"records": []interface{}{
				map[string]interface{}{
					"id":            "3477061000000419002",
					"Last_Name":     "Johnson",
					"First_Name":    "Emily",
					"Email":         "emily.johnson@example.com",
					"Company":       "XYZ Inc.",
					"Phone":         "+1-555-123-4567",
					"Lead_Source":   "Website",
					"Lead_Status":   "New",
					"Created_Time":  "2023-03-10T09:30:45+05:30",
					"Modified_Time": "2023-03-10T09:30:45+05:30",
				},
			},
			"count":  1,
			"module": "Leads",
		},
	}
}

func (t *RecordCreatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("record_created", "Record Created")

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

	schema := form.Build()

	return schema
}

func (t *RecordCreatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *RecordCreatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *RecordCreatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[recordCreatedTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	// Construct the search endpoint
	var endpoint string

	lr, err := ctx.GetMetadata("lastrun")

	if lr == nil {
		// If no last run time, fetch all records
		endpoint = fmt.Sprintf("https://www.zohoapis.com/crm/v7/%s/search", input.Module)
	} else {
		// If last run time exists, use it in the criteria
		lastRunTime := lr.(*time.Time)
		lastRunFormatted := lastRunTime.UTC().Format("2006-01-02T15:04:05+00:00")

		// URL encode the criteria to handle special characters
		encodedCriteria := url.QueryEscape(fmt.Sprintf("(Created_Time:greater_than:%s)", lastRunFormatted))

		endpoint = fmt.Sprintf("https://www.zohoapis.com/crm/v7/%s/search?criteria=%s",
			input.Module,
			encodedCriteria)
	}

	// Make API call
	result, err := shared.GetZohoCRMClient(token, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error calling Zoho CRM API: %v", err)
	}

	// Process response
	data, ok := result["data"].([]interface{})
	if !ok {
		if result["data"] == nil {
			return map[string]interface{}{
				"records": []interface{}{},
				"count":   0,
				"module":  input.Module,
			}, nil
		}
		return nil, errors.New("invalid response format: data field is not an array")
	}

	return map[string]interface{}{
		"records": data,
		"count":   len(data),
		"module":  input.Module,
	}, nil
}

func (t *RecordCreatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *RecordCreatedTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func (t *RecordCreatedTrigger) SampleData() sdkcore.JSON {
	return map[string]interface{}{
		"records": []interface{}{
			map[string]interface{}{
				"id":            "3477061000000419002",
				"Last_Name":     "Johnson",
				"First_Name":    "Emily",
				"Email":         "emily.johnson@example.com",
				"Company":       "XYZ Inc.",
				"Phone":         "+1-555-123-4567",
				"Lead_Source":   "Website",
				"Lead_Status":   "New",
				"Created_Time":  "2023-03-10T09:30:45+05:30",
				"Modified_Time": "2023-03-10T09:30:45+05:30",
			},
		},
		"count":  1,
		"module": "Leads",
	}
}

func NewRecordCreatedTrigger() sdk.Trigger {
	return &RecordCreatedTrigger{}
}
