package actions

import (
	"context"
	"errors"

	"github.com/wakflo/extensions/internal/integrations/googlesheets/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type addColumnInWorksheetActionProps struct {
	SpreadSheetID    string `json:"spreadsheetId"`
	SheetID          string `json:"sheetId"`
	SheetColumnIndex string `json:"sheetColumnIndex"`
}

type AddColumnInWorksheetAction struct{}

func (a *AddColumnInWorksheetAction) Name() string {
	return "Add Column in Worksheet"
}

func (a *AddColumnInWorksheetAction) Description() string {
	return "Add Column integration action allows you to dynamically add a new column to an existing worksheet within your workflow automation process. This action enables you to create custom fields or data points that can be used to store and manipulate information, further streamlining your workflow's efficiency and productivity."
}

func (a *AddColumnInWorksheetAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *AddColumnInWorksheetAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &addColumnInWorksheetDocs,
	}
}

func (a *AddColumnInWorksheetAction) Icon() *string {
	return nil
}

func (a *AddColumnInWorksheetAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"spreadsheetId": shared.GetSpreadsheetsInput("Spreadsheet", "spreadsheet ID", true),
		"sheetId":       shared.GetSheetIDInput("Sheet", "select sheet", true),
		"sheetColumnIndex": autoform.NewShortTextField().
			SetDisplayName("Sheet Column Index").
			SetDescription("The index of the column where the new column should be added. Index is zero based.").
			SetRequired(true).
			Build(),
	}
}

func (a *AddColumnInWorksheetAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[addColumnInWorksheetActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	sheetService, err := sheets.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
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

func (a *AddColumnInWorksheetAction) Auth() *sdk.Auth {
	return nil
}

func (a *AddColumnInWorksheetAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *AddColumnInWorksheetAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewAddColumnInWorksheetAction() sdk.Action {
	return &AddColumnInWorksheetAction{}
}
