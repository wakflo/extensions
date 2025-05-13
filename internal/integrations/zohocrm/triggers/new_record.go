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

type recordUpdatedTriggerProps struct {
	Module string `json:"module"`
}

type RecordUpdatedTrigger struct{}

func (t *RecordUpdatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "record_updated",
		DisplayName:   "Record Updated",
		Description:   "Triggers when an existing record is updated in a specified Zoho CRM module",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: newRecordDocs,
		SampleOutput: map[string]interface{}{
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
		},
	}
}

func (t *RecordUpdatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("new_record", "New Record")

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

func (t *RecordUpdatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *RecordUpdatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *RecordUpdatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[recordUpdatedTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	var endpoint string

	lr, err := ctx.GetMetadata("lastrun")

	if lr == nil {
		endpoint = fmt.Sprintf("https://www.zohoapis.com/crm/v7/%s/search", input.Module)
	} else {
		lastRunTime := lr.(*time.Time)
		lastRunFormatted := lastRunTime.UTC().Format("2006-01-02T15:04:05+00:00")

		encodedCriteria := url.QueryEscape(fmt.Sprintf("(Modified_Time:greater_than:%s)", lastRunFormatted))

		endpoint = fmt.Sprintf("https://www.zohoapis.com/crm/v7/%s/search?criteria=%s",
			input.Module,
			encodedCriteria)
	}

	result, err := shared.GetZohoCRMClient(ctx.Auth().AccessToken, http.MethodGet, endpoint, nil)
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

func (t *RecordUpdatedTrigger) Auth() *sdkcore.AuthMetadata {
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
