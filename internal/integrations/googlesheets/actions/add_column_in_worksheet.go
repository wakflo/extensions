package actions

import (
	"context"
	"errors"
	"strconv"

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
	ColumnName        string `json:"columnName"`
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

	form.TextField("columnName", "Column Name").
		Placeholder("Enter column name").
		HelpText("The name/header for the new column").
		Required(false)

	form.TextField("sheetColumnIndex", "Sheet Column Index").
		Placeholder("0").
		HelpText("The index of the column where the new column should be inserted BEFORE. Index is zero-based (0=A, 1=B, etc). For example, entering 5 will insert a new column before column F.").
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

	// Parse sheet ID
	sheetID, err := strconv.ParseInt(input.SheetID, 10, 64)
	if err != nil {
		return nil, errors.New("invalid sheet ID: must be a number")
	}

	// Parse column index
	columnIndex, err := strconv.ParseInt(input.SheetColumnIndex, 10, 64)
	if err != nil {
		return nil, errors.New("invalid column index: must be a number")
	}

	// Validate the column index is non-negative
	if columnIndex < 0 {
		return nil, errors.New("column index must be 0 or greater")
	}

	// Create a batch of requests
	requests := []*sheets.Request{}

	// First request: Insert the column
	insertDimensionRequest := &sheets.InsertDimensionRequest{
		Range: &sheets.DimensionRange{
			SheetId:    sheetID,
			Dimension:  "COLUMNS",
			StartIndex: columnIndex,
			EndIndex:   columnIndex + 1,
		},
		InheritFromBefore: false,
	}

	// Force send fields to ensure 0 values are sent
	insertDimensionRequest.Range.ForceSendFields = []string{"StartIndex", "EndIndex", "SheetId"}

	requests = append(requests, &sheets.Request{
		InsertDimension: insertDimensionRequest,
	})

	// Second request: Add column name if provided
	if input.ColumnName != "" {
		// Update the first cell of the new column with the column name
		updateCellRequest := &sheets.UpdateCellsRequest{
			Range: &sheets.GridRange{
				SheetId:          sheetID,
				StartRowIndex:    0,
				EndRowIndex:      1,
				StartColumnIndex: columnIndex,
				EndColumnIndex:   columnIndex + 1,
			},
			Rows: []*sheets.RowData{
				{
					Values: []*sheets.CellData{
						{
							UserEnteredValue: &sheets.ExtendedValue{
								StringValue: &input.ColumnName,
							},
						},
					},
				},
			},
			Fields: "userEnteredValue",
		}

		// Force send fields for the grid range
		updateCellRequest.Range.ForceSendFields = []string{"StartRowIndex", "EndRowIndex", "StartColumnIndex", "EndColumnIndex", "SheetId"}

		requests = append(requests, &sheets.Request{
			UpdateCells: updateCellRequest,
		})
	}

	// Create the BatchUpdateSpreadsheetRequest
	batchUpdateRequest := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}

	spreadsheet, err := sheetService.Spreadsheets.BatchUpdate(input.SpreadSheetID, batchUpdateRequest).
		Do()
	if err != nil {
		return nil, err
	}

	return spreadsheet, nil
}

func NewAddColumnInWorksheetAction() sdk.Action {
	return &AddColumnInWorksheetAction{}
}
