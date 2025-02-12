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

type findRecordActionProps struct {
	Bases        string `json:"bases"`
	Table        string `json:"table"`
	RecordID     string `json:"record-id"`
	SearchFields string `json:"search-fields,omitempty"`
	Views        string `json:"views"`
}

type FindRecordAction struct{}

func (a *FindRecordAction) Name() string {
	return "Find Record"
}

func (a *FindRecordAction) Description() string {
	return "Searches for a specific record within a database or data source and retrieves its details."
}

func (a *FindRecordAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *FindRecordAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &findRecordDocs,
	}
}

func (a *FindRecordAction) Icon() *string {
	return nil
}

func (a *FindRecordAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"record-id": autoform.NewShortTextField().
			SetDisplayName("Record ID").
			SetDescription("The ID of the record").
			SetRequired(true).Build(),
		"bases": shared.GetBasesInput(),
		"table": shared.GetTablesInput(),
	}
}

func (a *FindRecordAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[findRecordActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing airtable access token")
	}
	apiKey := ctx.Auth.Extra["api-key"]

	reqURL := fmt.Sprintf("%s/v0/%s/%s/%s", shared.BaseAPI, input.Bases, input.Table, input.RecordID)

	response, err := shared.AirtableRequest(apiKey, reqURL, http.MethodGet)
	if err != nil {
		return nil, errors.New("error fetching data")
	}

	return response, nil
}

func (a *FindRecordAction) Auth() *sdk.Auth {
	return nil
}

func (a *FindRecordAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *FindRecordAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewFindRecordAction() sdk.Action {
	return &FindRecordAction{}
}
