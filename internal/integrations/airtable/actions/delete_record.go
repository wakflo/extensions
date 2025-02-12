package actions

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/airtable/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type deleteRecordActionProps struct {
	RecordID string `json:"record-id"`
	Bases    string `json:"bases"`
	Table    string `json:"table"`
}

type DeleteRecordAction struct{}

func (a *DeleteRecordAction) Name() string {
	return "Delete Record"
}

func (a *DeleteRecordAction) Description() string {
	return "Deletes a record from a database or data source, permanently removing it from view and ensuring it is no longer accessible through the workflow."
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
		"record-id": autoform.NewShortTextField().
			SetDisplayName("Record ID").
			SetDescription("The ID of the record you want to delete. You can find the record ID by clicking on the record and then clicking on the share button. The ID will be in the URL.").
			SetRequired(true).Build(),
		"bases": shared.GetBasesInput(),
		"table": shared.GetTablesInput(),
	}
}

func (a *DeleteRecordAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[deleteRecordActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing airtable access token")
	}
	apiKey := ctx.Auth.Extra["api-key"]
	reqURL := fmt.Sprintf("%s/v0/%s/%s/%s", shared.BaseAPI, input.Bases, input.Table, input.RecordID)

	response, err := shared.AirtableRequest(apiKey, reqURL, http.MethodDelete)
	if err != nil {
		return nil, errors.New("error fetching data")
	}

	return response, nil
}

func (a *DeleteRecordAction) Auth() *sdk.Auth {
	return nil
}

func (a *DeleteRecordAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *DeleteRecordAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewDeleteRecordAction() sdk.Action {
	return &DeleteRecordAction{}
}
