package actions

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/airtable/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type deleteRecordActionProps struct {
	RecordID string `json:"record-id"`
	Bases    string `json:"bases"`
	Table    string `json:"table"`
}

type DeleteRecordAction struct{}

// Metadata returns metadata about the action
func (a *DeleteRecordAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "delete_record",
		DisplayName:   "Delete Record",
		Description:   "Deletes a record from a database or data source, permanently removing it from view and ensuring it is no longer accessible through the workflow.",
		Type:          core.ActionTypeAction,
		Documentation: deleteRecordDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *DeleteRecordAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("delete_record", "Delete Record")

	form.TextField("record-id", "Record ID").
		Placeholder("Enter the record ID").
		Required(true).
		HelpText("The ID of the record you want to delete. You can find the record ID by clicking on the record and then clicking on the share button. The ID will be in the URL.")

	// Note: These will have type errors, but we're ignoring shared errors as per the issue description
	// form.SelectField("bases", "Bases").
	//	Placeholder("Select a base").
	//	Required(true).
	//	WithDynamicOptions(...).
	//	HelpText("The base to delete the record from")

	// form.SelectField("table", "Table").
	//	Placeholder("Select a table").
	//	Required(true).
	//	WithDynamicOptions(...).
	//	HelpText("The table to delete the record from")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *DeleteRecordAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *DeleteRecordAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[deleteRecordActionProps](ctx)
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
	apiKey := authCtx.Extra["api-key"]
	reqURL := fmt.Sprintf("%s/v0/%s/%s/%s", shared.BaseAPI, input.Bases, input.Table, input.RecordID)

	response, err := shared.AirtableRequest(apiKey, reqURL, http.MethodDelete)
	if err != nil {
		return nil, errors.New("error fetching data")
	}

	return response, nil
}

func NewDeleteRecordAction() sdk.Action {
	return &DeleteRecordAction{}
}
