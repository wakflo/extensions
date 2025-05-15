package actions

import (
	"context"
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/googlesheets/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type getWorksheetByIdActionProps struct {
	SpreadSheetID     string `json:"spreadsheetId,omitempty"`
	SheetID           string `json:"sheetId"`
	IncludeTeamDrives bool   `json:"includeTeamDrives"`
}

type GetWorksheetByIdAction struct{}

// Metadata returns metadata about the action
func (a *GetWorksheetByIdAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_worksheet_by_id",
		DisplayName:   "Get Worksheet By ID",
		Description:   "Retrieves a specific worksheet by its unique identifier (ID), allowing you to access and manipulate its contents within your workflow.",
		Type:          core.ActionTypeAction,
		Documentation: getWorksheetByIdDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GetWorksheetByIdAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_worksheet_by_id", "Get Worksheet By ID")

	shared.RegisterSpreadsheetsProps(form, "spreadsheetId", "Spreadsheet", "spreadsheet ID", true)

	shared.RegisterSheetIDProps(form, true)

	form.CheckboxField("includeTeamDrives", "includeTeamDrives").
		Placeholder("Include Team Drives").
		HelpText("Include files from Team Drives in results")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetWorksheetByIdAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GetWorksheetByIdAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getWorksheetByIdActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	sheetService, err := sheets.NewService(context.Background(), option.WithTokenSource(*authCtx.TokenSource))
	if err != nil {
		return nil, err
	}

	if input.SpreadSheetID == "" {
		return nil, errors.New("spreadsheet ID is required")
	}

	if input.SheetID == "" {
		return nil, errors.New("sheet ID is required")
	}

	spreadsheet, err := sheetService.Spreadsheets.Get(input.SpreadSheetID).
		Do()
	if err != nil {
		return nil, err
	}

	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.SheetId == shared.ConvertToInt64(input.SheetID) {
			return sheet, nil
		}
	}
	return spreadsheet, err
}

func NewGetWorksheetByIdAction() sdk.Action {
	return &GetWorksheetByIdAction{}
}
