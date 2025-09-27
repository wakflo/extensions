package actions

import (
	"context"
	"errors"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/googlesheets/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type readRowInWorksheetActionProps struct {
	SpreadSheetID     string `json:"spreadsheetId,omitempty"`
	SheetID           string `json:"sheetId"`
	SheetTitle        string `json:"sheetTitle"`
	Range             string `json:"range"`
	SheetRow          string `json:"sheetRow"`
	IncludeTeamDrives bool   `json:"includeTeamDrives"`
}

type ReadRowInWorksheetAction struct{}

// Metadata returns metadata about the action
func (a *ReadRowInWorksheetAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "read_row_in_worksheet",
		DisplayName:   "Read Row in Worksheet",
		Description:   "Reads a single row from a worksheet and returns its values as an object. This action is useful when you need to retrieve specific data from a worksheet or perform actions based on the contents of a particular row.",
		Type:          core.ActionTypeAction,
		Documentation: readRowInWorksheetDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ReadRowInWorksheetAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("read_row_in_worksheet", "Read Row in Worksheet")

	shared.RegisterSpreadsheetsProps(form, "spreadsheetId", "Spreadsheet", "spreadsheet ID", true)

	shared.RegisterSheetTitleProps(form, true)

	form.TextField("sheetRow", "sheetRow").
		Placeholder("Sheet Row").
		HelpText("For adding data: use range format (e.g., A1:A100). For empty row: use row number (e.g., 5).").
		Required(true)

	form.CheckboxField("includeTeamDrives", "includeTeamDrives").
		Placeholder("Include Team Drives").
		HelpText("Include files from Team Drives in results")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ReadRowInWorksheetAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ReadRowInWorksheetAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[readRowInWorksheetActionProps](ctx)
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

	if input.SheetTitle == "" {
		return nil, errors.New("sheet title is required")
	}

	// Construct the range for reading the row
	readRange := fmt.Sprintf("%s!%s", input.SheetTitle, input.SheetRow)

	resp, err := sheetService.Spreadsheets.Values.Get(input.SpreadSheetID, readRange).Do()
	if err != nil {
		return nil, err
	}

	// Check if we got any values
	if len(resp.Values) == 0 {
		return nil, errors.New("no data found in the specified row")
	}

	rowValues := resp.Values[0]

	// Optionally, get the header row to create a key-value mapping
	// This assumes row 1 contains headers
	headerRange := fmt.Sprintf("%s!1:1", input.SheetTitle)
	headerResp, err := sheetService.Spreadsheets.Values.Get(input.SpreadSheetID, headerRange).Do()
	if err != nil {
		// If we can't get headers, just return the row values as an array
		return core.JSON(map[string]interface{}{
			"values": rowValues,
			"range":  resp.Range,
		}), nil
	}

	// Create a map with headers as keys if available
	result := make(map[string]interface{})
	if len(headerResp.Values) > 0 && len(headerResp.Values[0]) > 0 {
		headers := headerResp.Values[0]
		for i, header := range headers {
			if i < len(rowValues) {
				// Convert header to string and use as key
				key := fmt.Sprintf("%v", header)
				result[key] = rowValues[i]
			}
		}
	}

	return core.JSON(map[string]interface{}{
		"row":    result,
		"values": rowValues,
		"range":  resp.Range,
	}), nil
}

func NewReadRowInWorksheetAction() sdk.Action {
	return &ReadRowInWorksheetAction{}
}
