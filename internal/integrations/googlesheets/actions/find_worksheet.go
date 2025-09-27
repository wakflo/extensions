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

type findWorksheetActionProps struct {
	SpreadSheetID string `json:"spreadsheetId,omitempty"`
	SheetTitle    string `json:"sheetTitle"`
}

type FindWorksheetAction struct{}

// Metadata returns metadata about the action
func (a *FindWorksheetAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "find_worksheet",
		DisplayName:   "Find Worksheet",
		Description:   "Locates and retrieves a specific worksheet from a spreadsheet application, allowing you to automate tasks that rely on this worksheet's data.",
		Type:          core.ActionTypeAction,
		Documentation: findWorksheetDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"found":     true,
			"sheetId":   123456789,
			"title":     "Sheet1",
			"index":     0,
			"sheetType": "GRID",
			"gridProperties": map[string]any{
				"rowCount":    1000,
				"columnCount": 26,
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *FindWorksheetAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("find_worksheet", "Find Worksheet")

	shared.RegisterSpreadsheetsProps(form, "spreadsheetId", "Spreadsheet", "spreadsheet ID", true)

	shared.RegisterSheetTitleProps(form, true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *FindWorksheetAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *FindWorksheetAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[findWorksheetActionProps](ctx)
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

	spreadsheet, err := sheetService.Spreadsheets.Get(input.SpreadSheetID).Do()
	if err != nil {
		return nil, err
	}

	// Search for the sheet with matching title
	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.Title == input.SheetTitle {
			// Found the sheet, return its properties
			result := map[string]interface{}{
				"found":   true,
				"sheetId": sheet.Properties.SheetId,
				"title":   sheet.Properties.Title,
				"index":   sheet.Properties.Index,
			}

			// Add sheet type if available
			if sheet.Properties.SheetType != "" {
				result["sheetType"] = sheet.Properties.SheetType
			}

			// Add grid properties if available
			if sheet.Properties.GridProperties != nil {
				gridProps := map[string]interface{}{
					"rowCount":    sheet.Properties.GridProperties.RowCount,
					"columnCount": sheet.Properties.GridProperties.ColumnCount,
				}

				if sheet.Properties.GridProperties.FrozenRowCount > 0 {
					gridProps["frozenRowCount"] = sheet.Properties.GridProperties.FrozenRowCount
				}
				if sheet.Properties.GridProperties.FrozenColumnCount > 0 {
					gridProps["frozenColumnCount"] = sheet.Properties.GridProperties.FrozenColumnCount
				}

				result["gridProperties"] = gridProps
			}

			// Add tab color if available
			if sheet.Properties.TabColor != nil {
				result["tabColor"] = map[string]interface{}{
					"red":   sheet.Properties.TabColor.Red,
					"green": sheet.Properties.TabColor.Green,
					"blue":  sheet.Properties.TabColor.Blue,
					"alpha": sheet.Properties.TabColor.Alpha,
				}
			}

			return core.JSON(result), nil
		}
	}

	// Sheet not found - return proper JSON response
	return core.JSON(map[string]interface{}{
		"found":   false,
		"message": fmt.Sprintf("No sheet found with title '%s' in spreadsheet", input.SheetTitle),
	}), nil
}

func NewFindWorksheetAction() sdk.Action {
	return &FindWorksheetAction{}
}
