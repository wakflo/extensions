package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/airtable/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type updateRecordActionProps struct {
	Name  string `json:"name"`
	Bases string `json:"bases"`
	Table string `json:"table"`
}

type UpdateRecordAction struct{}

func (a *UpdateRecordAction) Name() string {
	return "Update Record"
}

func (a *UpdateRecordAction) Description() string {
	return "Updates an existing record in your database or application with new information, ensuring accurate and up-to-date data."
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
		"record-id": autoform.NewShortTextField().
			SetDisplayName("Record ID").
			SetDescription("The record's ID").
			SetRequired(true).Build(),
		"bases":  shared.GetBasesInput(),
		"table":  shared.GetTablesInput(),
		"fields": shared.GetFieldsInput(),
		"views":  shared.GetViewsInput(),
	}
}

func (a *UpdateRecordAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateRecordActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing airtable access token")
	}
	accessToken := ctx.Auth.Extra["api-key"]

	reqURL := fmt.Sprintf("%s/v0/meta/bases/%s/tables", shared.BaseAPI, input.Bases)

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response interface{}
	if errs := json.Unmarshal(body, &response); errs != nil {
		return nil, errors.New("error parsing response")
	}

	return response, nil
}

func (a *UpdateRecordAction) Auth() *sdk.Auth {
	return nil
}

func (a *UpdateRecordAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *UpdateRecordAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewUpdateRecordAction() sdk.Action {
	return &UpdateRecordAction{}
}
