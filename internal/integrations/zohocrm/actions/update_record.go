package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/zohocrm/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type updateRecordActionProps struct {
	Module   string `json:"module"`
	RecordID string `json:"recordId"`
	Data     string `json:"data"`
}

type UpdateRecordAction struct{}

// func (a *UpdateRecordAction) Name() string {
// 	return "Update Record"
// }

func (a *UpdateRecordAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "update_record",
		DisplayName:   "Update Record",
		Description:   "Updates an existing record in a specified Zoho CRM module with the provided data",
		Type:          sdkcore.ActionTypeAction,
		Documentation: updateRecordDocs,
		SampleOutput: map[string]interface{}{
			"id":            "3477061000000419001",
			"Modified_Time": "2023-01-15T16:30:45+05:30",
			"Modified_By": map[string]interface{}{
				"id":   "3477061000000319001",
				"name": "John Doe",
			},
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *UpdateRecordAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("update_record", "Update Record")

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

	form.TextField("recordId", "Record ID").
		Placeholder("Enter the record ID").
		Required(true).
		HelpText("Enter the ID of the record you want to update.")

		// form.CodeEditorField("data", "Record Data").
	// Placeholder("Enter JSON data").
	// Required(true)

	schema := form.Build()

	return schema
}

func (a *UpdateRecordAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateRecordActionProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

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
	result, err := shared.GetZohoCRMClient(token, http.MethodPut, endpoint, requestData)
	if err != nil {
		return nil, fmt.Errorf("error calling Zoho CRM API: %v", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok || len(data) == 0 {
		return nil, errors.New("invalid response format: data field is missing or empty")
	}

	return data[0], nil
}

func (a *UpdateRecordAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewUpdateRecordAction() sdk.Action {
	return &UpdateRecordAction{}
}
