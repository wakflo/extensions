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

type recordUpdatedTriggerProps struct {
	Module string `json:"module"`
}

type RecordUpdatedTrigger struct{}

func (t *RecordUpdatedTrigger) Name() string {
	return "Record Updated"
}

func (t *RecordUpdatedTrigger) Description() string {
	return "Triggers when an existing record is updated in a specified Zoho CRM module"
}

func (t *RecordUpdatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *RecordUpdatedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &newRecordDocs,
	}
}

func (t *RecordUpdatedTrigger) Icon() *string {
	return nil
}

func (t *RecordUpdatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"module": autoform.NewDynamicField(sdkcore.String).
			SetDisplayName("Module").
			SetDescription("The Zoho CRM module to monitor for newly created records (e.g., Leads, Contacts, Accounts)").
			SetDynamicOptions(shared.GetModulesFunction()).
			SetRequired(true).
			Build(),
	}
}

func (t *RecordUpdatedTrigger) Start(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *RecordUpdatedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *RecordUpdatedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[recordUpdatedTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	var endpoint string
	if ctx.Metadata().LastRun == nil {
		endpoint = fmt.Sprintf("https://www.zohoapis.com/crm/v7/%s/search", input.Module)
	} else {
		lastRunTime := *ctx.Metadata().LastRun
		lastRunFormatted := lastRunTime.UTC().Format("2006-01-02T15:04:05+00:00")

		encodedCriteria := url.QueryEscape(fmt.Sprintf("(Modified_Time:greater_than:%s)", lastRunFormatted))

		endpoint = fmt.Sprintf("https://www.zohoapis.com/crm/v7/%s/search?criteria=%s",
			input.Module,
			encodedCriteria)
	}

	result, err := shared.GetZohoCRMClient(ctx.Auth.AccessToken, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error calling Zoho CRM API: %v", err)
	}

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

func (t *RecordUpdatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *RecordUpdatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *RecordUpdatedTrigger) SampleData() sdkcore.JSON {
	return map[string]interface{}{
		"records": []interface{}{
			map[string]interface{}{
				"id":            "3477061000000419001",
				"Last_Name":     "Smith",
				"First_Name":    "John",
				"Email":         "john.smith@example.com",
				"Created_Time":  "2023-01-15T15:45:30+05:30",
				"Modified_Time": "2023-01-16T10:20:15+05:30",
			},
		},
		"count":  1,
		"module": "Leads",
	}
}

func NewRecordUpdatedTrigger() sdk.Trigger {
	return &RecordUpdatedTrigger{}
}
