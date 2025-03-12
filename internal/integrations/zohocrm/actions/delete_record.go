package actions

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/zohocrm/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type deleteRecordActionProps struct {
	Module   string `json:"module"`
	RecordID string `json:"recordId"`
}

type DeleteRecordAction struct{}

func (a *DeleteRecordAction) Name() string {
	return "Delete Record"
}

func (a *DeleteRecordAction) Description() string {
	return "Deletes a specific record from a Zoho CRM module"
}

func (a *DeleteRecordAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *DeleteRecordAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &deleteRecordDocs,
	}
}

func (a *DeleteRecordAction) Icon() *string {
	return nil
}

func (a *DeleteRecordAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"module": autoform.NewDynamicField(sdkcore.String).
			SetDisplayName("Module").
			SetDescription("The Zoho CRM module from which to delete the record (e.g., Leads, Contacts, Accounts)").
			SetDynamicOptions(shared.GetModulesFunction()).
			SetRequired(true).
			Build(),
		"recordId": autoform.NewShortTextField().
			SetDisplayName("Record ID").
			SetDescription("The ID of the record to delete").
			SetRequired(true).
			Build(),
	}
}

func (a *DeleteRecordAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[deleteRecordActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/%s", input.Module, input.RecordID)
	result, err := shared.GetZohoCRMClient(ctx.Auth.AccessToken, http.MethodDelete, endpoint, nil)
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

func (a *DeleteRecordAction) Auth() *sdk.Auth {
	return nil
}

func (a *DeleteRecordAction) SampleData() sdkcore.JSON {
	return map[string]interface{}{
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
	}
}

func (a *DeleteRecordAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewDeleteRecordAction() sdk.Action {
	return &DeleteRecordAction{}
}
