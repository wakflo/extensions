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

type addColumnInWorksheetActionProps struct {
	SpreadSheetID     string `json:"spreadsheetId"`
	SheetID           string `json:"sheetId"`
	SheetColumnIndex  string `json:"sheetColumnIndex"`
	IncludeTeamDrives bool   `json:"includeTeamDrives"`
}

type AddColumnInWorksheetAction struct{}

// Metadata returns metadata about the action
func (a *AddColumnInWorksheetAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "add_column_in_worksheet",
		DisplayName:   "Add Column in Worksheet",
		Description:   "Add Column integration action allows you to dynamically add a new column to an existing worksheet within your workflow automation process. This action enables you to create custom fields or data points that can be used to store and manipulate information, further streamlining your workflow's efficiency and productivity.",
		Type:          core.ActionTypeAction,
		Documentation: addColumnInWorksheetDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *AddColumnInWorksheetAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("add_column_in_worksheet", "Add Column in Worksheet")

	shared.RegisterSpreadsheetsProps(form, "spreadsheetId", "SpreadSheet", "Select Spreadsheet", true)

	shared.RegisterSheetIDProps(form, true)

	form.CheckboxField("includeTeamDrives", "Include Team Drives").
		Placeholder("Enter a value for Include Team Drives.").
		Required(false).
		HelpText("Whether to include team drives in the folder selection.")

	form.TextField("sheetColumnIndex", "sheetColumnIndex").
		Placeholder("Sheet Column Index").
		HelpText("The index of the column where the new column should be added. Index is zero based.(E.g; 0,1 etc)").
		Required(true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *AddColumnInWorksheetAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *AddColumnInWorksheetAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[addColumnInWorksheetActionProps](ctx)
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
	if input.SheetColumnIndex == "" {
		return nil, errors.New("sheet column index is required")
	}

	// Create the InsertDimensionRequest
	insertDimensionRequest := &sheets.InsertDimensionRequest{
		Range: &sheets.DimensionRange{
			SheetId:    shared.ConvertToInt64(input.SheetID),
			Dimension:  "COLUMNS",
			StartIndex: shared.ConvertToInt64(input.SheetColumnIndex),
			EndIndex:   shared.ConvertToInt64(input.SheetColumnIndex) + 1,
		},
		InheritFromBefore: false,
	}

	// Create the BatchUpdateSpreadsheetRequest
	batchUpdateRequest := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				InsertDimension: insertDimensionRequest,
			},
		},
	}

	spreadsheet, err := sheetService.Spreadsheets.BatchUpdate(input.SpreadSheetID, batchUpdateRequest).
		Do()

	return spreadsheet, err
}

func NewAddColumnInWorksheetAction() sdk.Action {
	return &AddColumnInWorksheetAction{}
}
