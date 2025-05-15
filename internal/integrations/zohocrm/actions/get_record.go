package actions

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/zohocrm/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type getRecordActionProps struct {
	Module   string `json:"module"`
	RecordID string `json:"recordId"`
}

type GetRecordAction struct{}

func (a *GetRecordAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_record",
		DisplayName:   "Get Record",
		Description:   "Retrieves a specific record from a Zoho CRM module using its ID",
		Type:          sdkcore.ActionTypeAction,
		Documentation: getRecordDocs,
		SampleOutput: map[string]interface{}{
			"id":            "3477061000000419001",
			"Last_Name":     "Smith",
			"First_Name":    "John",
			"Email":         "john.smith@example.com",
			"Company":       "ACME Corp",
			"Lead_Status":   "Qualified",
			"Created_Time":  "2023-01-15T15:45:30+05:30",
			"Modified_Time": "2023-01-15T15:45:30+05:30",
			"Created_By": map[string]interface{}{
				"id":   "3477061000000319001",
				"name": "Jane Doe",
			},
			"Modified_By": map[string]interface{}{
				"id":   "3477061000000319001",
				"name": "Jane Doe",
			},
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *GetRecordAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_record", "Get Record")

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

	form.TextField("record_id", "Record ID").
		Placeholder("Enter the record ID").
		Required(true).
		HelpText("The unique identifier of the record to be retrieved.")

	schema := form.Build()

	return schema
}

func (a *GetRecordAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getRecordActionProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	endpoint := fmt.Sprintf("%s/%s", input.Module, input.RecordID)
	result, err := shared.GetZohoCRMClient(token, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error calling Zoho CRM API: %v", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok || len(data) == 0 {
		return nil, errors.New("invalid response format: data field is missing or empty")
	}

	return data[0], nil
}

func (a *GetRecordAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewGetRecordAction() sdk.Action {
	return &GetRecordAction{}
}
