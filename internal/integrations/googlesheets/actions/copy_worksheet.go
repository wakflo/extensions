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

type copyWorksheetActionProps struct {
	SourceSpreadSheetID      string `json:"spreadsheetId,omitempty"`
	SourceSheetTitle         string `json:"sheetTitle"`
	DestinationSpreadSheetID string `json:"destinationSpreadsheetId,omitempty"`
	NewSheetTitle            string `json:"newSheetTitle,omitempty"`
	IncludeTeamDrives        bool   `json:"includeTeamDrives"`
}

type CopyWorksheetAction struct{}

// Metadata returns metadata about the action
func (a *CopyWorksheetAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "copy_worksheet",
		DisplayName:   "Copy Worksheet",
		Description:   "Copies a worksheet within the same spreadsheet or to a different spreadsheet. The copy includes all data, formatting, and formulas from the source worksheet.",
		Type:          core.ActionTypeAction,
		Documentation: copyWorksheetDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"sheetId":                  123456,
			"title":                    "Copy of Sheet1",
			"index":                    2,
			"destinationSpreadsheetId": "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CopyWorksheetAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("copy_worksheet", "Copy Worksheet")

	shared.RegisterSpreadsheetsProps(form, "spreadsheetId", "Source Spreadsheet", "source spreadsheet ID", true)

	shared.RegisterSheetTitleProps(form, true)

	shared.RegisterSpreadsheetsProps(form, "destinationSpreadsheetId", "Destination Spreadsheet", "destination spreadsheet ID (leave empty to copy within same spreadsheet)", false)

	form.TextField("newSheetTitle", "newSheetTitle").
		Placeholder("New Sheet Title").
		HelpText("Title for the copied worksheet. If empty, will use 'Copy of [original name]'")

	form.CheckboxField("includeTeamDrives", "includeTeamDrives").
		Placeholder("Include Team Drives").
		HelpText("Include files from Team Drives in results")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CopyWorksheetAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CopyWorksheetAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[copyWorksheetActionProps](ctx)
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

	if input.SourceSpreadSheetID == "" {
		return nil, errors.New("source spreadsheet ID is required")
	}

	if input.SourceSheetTitle == "" {
		return nil, errors.New("source sheet title is required")
	}

	// Get the spreadsheet to find the sheet ID by title
	spreadsheet, err := sheetService.Spreadsheets.Get(input.SourceSpreadSheetID).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get source spreadsheet: %w", err)
	}

	var sourceSheetID int64
	found := false
	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.Title == input.SourceSheetTitle {
			sourceSheetID = sheet.Properties.SheetId
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("sheet with title '%s' not found in source spreadsheet", input.SourceSheetTitle)
	}

	// Determine destination spreadsheet ID
	destinationSpreadsheetID := input.DestinationSpreadSheetID
	if destinationSpreadsheetID == "" {
		// Copy within the same spreadsheet
		destinationSpreadsheetID = input.SourceSpreadSheetID
	}

	// Create the copy sheet request
	copyRequest := &sheets.CopySheetToAnotherSpreadsheetRequest{
		DestinationSpreadsheetId: destinationSpreadsheetID,
	}

	// Execute the copy operation
	resp, err := sheetService.Spreadsheets.Sheets.CopyTo(
		input.SourceSpreadSheetID,
		sourceSheetID,
		copyRequest,
	).Do()
	if err != nil {
		return nil, err
	}

	// If a new title was specified, update it
	if input.NewSheetTitle != "" {
		updateRequest := &sheets.BatchUpdateSpreadsheetRequest{
			Requests: []*sheets.Request{
				{
					UpdateSheetProperties: &sheets.UpdateSheetPropertiesRequest{
						Properties: &sheets.SheetProperties{
							SheetId: resp.SheetId,
							Title:   input.NewSheetTitle,
						},
						Fields: "title",
					},
				},
			},
		}

		_, err = sheetService.Spreadsheets.BatchUpdate(
			destinationSpreadsheetID,
			updateRequest,
		).Do()

		if err != nil {
			fmt.Printf("Warning: Failed to rename copied sheet: %v\n", err)
		} else {
			resp.Title = input.NewSheetTitle
		}
	}

	return core.JSON(map[string]interface{}{
		"sheetId":                  resp.SheetId,
		"title":                    resp.Title,
		"index":                    resp.Index,
		"destinationSpreadsheetId": destinationSpreadsheetID,
		"success":                  true,
	}), nil
}

func NewCopyWorksheetAction() sdk.Action {
	return &CopyWorksheetAction{}
}
