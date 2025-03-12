package triggers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/wakflo/extensions/internal/integrations/zohocrm/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type recordCreatedTriggerProps struct {
	Module string `json:"module"`
}

type RecordCreatedTrigger struct{}

func (t *RecordCreatedTrigger) Name() string {
	return "Record Created"
}

func (t *RecordCreatedTrigger) Description() string {
	return "Triggers when a new record is created in a specified Zoho CRM module"
}

func (t *RecordCreatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *RecordCreatedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &recordCreatedDocs,
	}
}

func (t *RecordCreatedTrigger) Icon() *string {
	return nil
}

func (t *RecordCreatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"module": autoform.NewDynamicField(sdkcore.String).
			SetDisplayName("Module").
			SetDescription("The Zoho CRM module to monitor for newly created records (e.g., Leads, Contacts, Accounts)").
			SetDynamicOptions(shared.GetModulesFunction()).
			SetRequired(true).
			Build(),
	}
}

func (t *RecordCreatedTrigger) Start(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *RecordCreatedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *RecordCreatedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[recordCreatedTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	// Construct the search endpoint
	var endpoint string
	if ctx.Metadata().LastRun == nil {
		// If no last run time, fetch all records
		endpoint = fmt.Sprintf("https://www.zohoapis.com/crm/v7/%s/search", input.Module)
	} else {
		// If last run time exists, use it in the criteria
		lastRunTime := *ctx.Metadata().LastRun
		lastRunFormatted := lastRunTime.UTC().Format("2006-01-02T15:04:05+00:00")

		// URL encode the criteria to handle special characters
		encodedCriteria := url.QueryEscape(fmt.Sprintf("(Created_Time:greater_than:%s)", lastRunFormatted))

		endpoint = fmt.Sprintf("https://www.zohoapis.com/crm/v7/%s/search?criteria=%s",
			input.Module,
			encodedCriteria)
	}

	// Make API call
	result, err := shared.GetZohoCRMClient(ctx.Auth.AccessToken, http.MethodGet, endpoint, nil)
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

func (t *RecordCreatedTrigger) Auth() *sdk.Auth {
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
