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

type findRecordActionProps struct {
	Bases        string `json:"bases"`
	Table        string `json:"table"`
	RecordID     string `json:"record-id"`
	SearchFields string `json:"search-fields,omitempty"`
	Views        string `json:"views"`
}

type FindRecordAction struct{}

// Metadata returns metadata about the action
func (a *FindRecordAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "find_record",
		DisplayName:   "Find Record",
		Description:   "Searches for a specific record within a database or data source and retrieves its details.",
		Type:          core.ActionTypeAction,
		Documentation: findRecordDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *FindRecordAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("find_record", "Find Record")

	form.TextField("record-id", "Record ID").
		Placeholder("Enter the record ID").
		Required(true).
		HelpText("The ID of the record")

	shared.RegisterBasesProps(form)

	shared.RegisterTablesProps(form)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *FindRecordAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *FindRecordAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[findRecordActionProps](ctx)
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

	response, err := shared.AirtableRequest(apiKey, reqURL, http.MethodGet)
	if err != nil {
		return nil, errors.New("error fetching data")
	}

	return response, nil
}

func NewFindRecordAction() sdk.Action {
	return &FindRecordAction{}
}
