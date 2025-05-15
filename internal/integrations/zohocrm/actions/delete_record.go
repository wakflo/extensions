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

type deleteRecordActionProps struct {
	Module   string `json:"module"`
	RecordID string `json:"recordId"`
}

type DeleteRecordAction struct{}

func (a *DeleteRecordAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "delete_record",
		DisplayName:   "Delete Record",
		Description:   "Deletes a specific record from a Zoho CRM module",
		Type:          sdkcore.ActionTypeAction,
		Documentation: deleteRecordDocs,
		SampleOutput: map[string]interface{}{
			"success": true,
			"id":      "3477061000000419001",
			"detail": map[string]interface{}{
				"code":    "SUCCESS",
				"status":  "success",
				"message": "record deleted",
				"details": map[string]interface{}{
					"id": "3477061000000419001",
				},
			},
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *DeleteRecordAction) Properties() *smartform.FormSchema {

	form := smartform.NewForm("delete_record", "Delete Record")

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
		HelpText("The unique identifier of the record to be deleted.")

	schema := form.Build()

	return schema
}

func (a *DeleteRecordAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[deleteRecordActionProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	endpoint := fmt.Sprintf("%s/%s", input.Module, input.RecordID)
	result, err := shared.GetZohoCRMClient(token, http.MethodDelete, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error calling Zoho CRM API: %v", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok || len(data) == 0 {
		return nil, errors.New("invalid response format: data field is missing or empty")
	}

	return map[string]interface{}{
		"success": true,
		"id":      input.RecordID,
		"detail":  data[0],
	}, nil
}

func (a *DeleteRecordAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewDeleteRecordAction() sdk.Action {
	return &DeleteRecordAction{}
}
