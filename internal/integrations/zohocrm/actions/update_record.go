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

type updateRecordActionProps struct {
	Module   string `json:"module"`
	RecordID string `json:"recordId"`
	Data     string `json:"data"`
}

type UpdateRecordAction struct{}

func (a *UpdateRecordAction) Name() string {
	return "Update Record"
}

func (a *UpdateRecordAction) Description() string {
	return "Updates an existing record in a specified Zoho CRM module with the provided data"
}

func (a *UpdateRecordAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *UpdateRecordAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &updateRecordDocs,
	}
}

func (a *UpdateRecordAction) Icon() *string {
	return nil
}

func (a *UpdateRecordAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"module": autoform.NewDynamicField(sdkcore.String).
			SetDisplayName("Module").
			SetDescription("The Zoho CRM module where the record will be updated (e.g., Leads, Contacts, Accounts)").
			SetDynamicOptions(shared.GetModulesFunction()).
			SetRequired(true).
			Build(),
		"recordId": autoform.NewShortTextField().
			SetDisplayName("Record ID").
			SetDescription("The ID of the record to update").
			SetRequired(true).
			Build(),
		"data": autoform.NewCodeEditorField().
			SetDisplayName("Record Data").
			SetDescription("The updated data for the record in JSON format").
			SetRequired(true).
			Build(),
	}
}

func (a *UpdateRecordAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateRecordActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	trimmedData := strings.TrimSpace(input.Data)

	if !strings.HasPrefix(trimmedData, "{") || !strings.HasSuffix(trimmedData, "}") {
		return nil, errors.New("invalid JSON format: data must be a valid JSON object starting with '{' and ending with '}'")
	}

	var dataMap map[string]interface{}
	if err := json.Unmarshal([]byte(trimmedData), &dataMap); err != nil {
		return nil, fmt.Errorf("invalid JSON data: %v. Please ensure you're providing a valid JSON object", err)
	}

	// Check if the data map is empty
	if len(dataMap) == 0 {
		return nil, errors.New("empty data: please provide at least one field to update")
	}

	requestData := map[string]interface{}{
		"data": []interface{}{dataMap},
	}

	endpoint := fmt.Sprintf("%s/%s", input.Module, input.RecordID)
	result, err := shared.GetZohoCRMClient(ctx.Auth.AccessToken, http.MethodPut, endpoint, requestData)
	if err != nil {
		return nil, fmt.Errorf("error calling Zoho CRM API: %v", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok || len(data) == 0 {
		return nil, errors.New("invalid response format: data field is missing or empty")
	}

	return data[0], nil
}

func (a *UpdateRecordAction) Auth() *sdk.Auth {
	return nil
}

func (a *UpdateRecordAction) SampleData() sdkcore.JSON {
	return map[string]interface{}{
		"id":            "3477061000000419001",
		"Modified_Time": "2023-01-15T16:30:45+05:30",
		"Modified_By": map[string]interface{}{
			"id":   "3477061000000319001",
			"name": "John Doe",
		},
	}
}

func (a *UpdateRecordAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewUpdateRecordAction() sdk.Action {
	return &UpdateRecordAction{}
}
