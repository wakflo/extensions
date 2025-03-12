package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/wakflo/extensions/internal/integrations/zohocrm/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createRecordActionProps struct {
	Module string `json:"module"`
	Data   string `json:"data"`
}

type CreateRecordAction struct{}

func (a *CreateRecordAction) Name() string {
	return "Create Record"
}

func (a *CreateRecordAction) Description() string {
	return "Creates a new record in a specified Zoho CRM module with the provided data"
}

func (a *CreateRecordAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateRecordAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createRecordDocs,
	}
}

func (a *CreateRecordAction) Icon() *string {
	return nil
}

func (a *CreateRecordAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"module": autoform.NewDynamicField(sdkcore.String).
			SetDisplayName("Module").
			SetDescription("The Zoho CRM module where the record will be created (e.g., Leads, Contacts, Accounts)").
			SetDynamicOptions(shared.GetModulesFunction()).
			SetRequired(true).
			Build(),
		"data": autoform.NewCodeEditorField().
			SetDisplayName("Record Data").
			SetDescription("The data for the new record in JSON format").
			SetRequired(true).
			Build(),
	}
}

func (a *CreateRecordAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createRecordActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	var dataMap map[string]interface{}
	if err := json.Unmarshal([]byte(input.Data), &dataMap); err != nil {
		return nil, fmt.Errorf("invalid JSON data: %v", err)
	}

	requestData := map[string]interface{}{
		"data": []interface{}{dataMap},
	}

	endpoint := input.Module
	result, err := shared.GetZohoCRMClient(ctx.Auth.AccessToken, http.MethodPost, endpoint, requestData)
	if err != nil {
		if strings.Contains(err.Error(), "201") {
			return map[string]interface{}{
				"success": true,
				"message": "Record created successfully",
				"module":  input.Module,
				"data":    dataMap,
			}, nil
		}
		return nil, fmt.Errorf("error calling Zoho CRM API: %v", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok || len(data) == 0 {
		return nil, errors.New("invalid response format: data field is missing or empty")
	}

	return data[0], nil
}

func (a *CreateRecordAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateRecordAction) SampleData() sdkcore.JSON {
	return map[string]interface{}{
		"id":            "3477061000000419001",
		"Created_Time":  "2023-01-15T15:45:30+05:30",
		"Modified_Time": "2023-01-15T15:45:30+05:30",
		"Created_By": map[string]interface{}{
			"id":   "3477061000000319001",
			"name": "John Doe",
		},
		"Modified_By": map[string]interface{}{
			"id":   "3477061000000319001",
			"name": "John Doe",
		},
	}
}

func (a *CreateRecordAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateRecordAction() sdk.Action {
	return &CreateRecordAction{}
}
