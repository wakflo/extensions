package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/airtable/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type updateRecordActionProps struct {
	Name     string `json:"name"`
	Bases    string `json:"bases"`
	Table    string `json:"table"`
	RecordID string `json:"record-id"`
}

type UpdateRecordAction struct{}

// Metadata returns metadata about the action
func (a *UpdateRecordAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "update_record",
		DisplayName:   "Update Record",
		Description:   "Updates an existing record in your database or application with new information, ensuring accurate and up-to-date data.",
		Type:          core.ActionTypeAction,
		Documentation: updateRecordDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *UpdateRecordAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("update_record", "Update Record")

	form.TextField("record-id", "Record ID").
		Placeholder("Enter the record ID").
		Required(true).
		HelpText("The record's ID")

	shared.RegisterBasesProps(form)

	shared.RegisterTablesProps(form)

	shared.RegisterFieldsProps(form)

	shared.RegisterViewsProps(form)

	// form.SelectField("views", "Views").
	//	Placeholder("Select views").
	//	Required(false).
	//	WithDynamicOptions(...).
	//	HelpText("The views to update")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *UpdateRecordAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *UpdateRecordAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[updateRecordActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["api-key"] == "" {
		return nil, errors.New("missing airtable access token")
	}
	accessToken := authCtx.Extra["api-key"]

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

func NewUpdateRecordAction() sdk.Action {
	return &UpdateRecordAction{}
}
