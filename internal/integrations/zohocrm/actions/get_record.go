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

type getRecordActionProps struct {
	Module   string `json:"module"`
	RecordID string `json:"recordId"`
}

type GetRecordAction struct{}

func (a *GetRecordAction) Name() string {
	return "Get Record"
}

func (a *GetRecordAction) Description() string {
	return "Retrieves a specific record from a Zoho CRM module using its ID"
}

func (a *GetRecordAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetRecordAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getRecordDocs,
	}
}

func (a *GetRecordAction) Icon() *string {
	return nil
}

func (a *GetRecordAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"module": autoform.NewDynamicField(sdkcore.String).
			SetDisplayName("Module").
			SetDescription("The Zoho CRM module to retrieve the record from (e.g., Leads, Contacts, Accounts)").
			SetDynamicOptions(shared.GetModulesFunction()).
			SetRequired(true).
			Build(),
		"recordId": autoform.NewShortTextField().
			SetDisplayName("Record ID").
			SetDescription("The ID of the record to retrieve").
			SetRequired(true).
			Build(),
	}
}

func (a *GetRecordAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getRecordActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/%s", input.Module, input.RecordID)
	result, err := shared.GetZohoCRMClient(ctx.Auth.AccessToken, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error calling Zoho CRM API: %v", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok || len(data) == 0 {
		return nil, errors.New("invalid response format: data field is missing or empty")
	}

	return data[0], nil
}

func (a *GetRecordAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetRecordAction) SampleData() sdkcore.JSON {
	return map[string]interface{}{
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
	}
}

func (a *GetRecordAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetRecordAction() sdk.Action {
	return &GetRecordAction{}
}
