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

type createRecordActionProps struct {
	Module string `json:"module"`
	Data   string `json:"data"`
}

type CreateRecordAction struct{}

func (a *CreateRecordAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_record",
		DisplayName:   "Create Record",
		Description:   "Creates a new record in a specified Zoho CRM module with the provided data",
		Type:          sdkcore.ActionTypeAction,
		Documentation: createRecordDocs,
		SampleOutput: map[string]any{
			"message": "hello world",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *CreateRecordAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_record", "Create Record")

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

	// form.CodeEditorField("data", "Record Data").
	// Placeholder("Enter JSON data").
	// Required(true)

	schema := form.Build()

	return schema

}

func (a *CreateRecordAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createRecordActionProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	var dataMap map[string]interface{}
	if err := json.Unmarshal([]byte(input.Data), &dataMap); err != nil {
		return nil, fmt.Errorf("invalid JSON data: %v", err)
	}

	requestData := map[string]interface{}{
		"data": []interface{}{dataMap},
	}

	endpoint := input.Module
	result, err := shared.GetZohoCRMClient(token, http.MethodPost, endpoint, requestData)
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

func (a *CreateRecordAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewCreateRecordAction() sdk.Action {
	return &CreateRecordAction{}
}
