package googlesheets

import (
	"context"
	"errors"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type addColumnInWorkSheetOperationProps struct {
	SpreadSheetID    string `json:"spreadSheetId"`
	SheetID          string `json:"sheetId"`
	SheetColumnIndex string `json:"sheetColumnIndex"`
}

type AddColumnInWorkSheetOperation struct {
	options *sdk.OperationInfo
}

func NewAddColumnInWorkSheetOperation() *AddColumnInWorkSheetOperation {
	return &AddColumnInWorkSheetOperation{
		options: &sdk.OperationInfo{
			Name:        "Add Column in Sheet",
			Description: "Add a column in a sheet",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"spreadSheetId": autoform.NewShortTextField().
					SetDisplayName("Spreadsheet ID").
					SetDescription("The ID of the spreadsheet.").
					SetRequired(true).
					Build(),
				"sheetId": autoform.NewShortTextField().
					SetDisplayName("Sheet ID").
					SetDescription("The ID of the sheet.").
					SetRequired(true).
					Build(),
				"sheetColumnIndex": autoform.NewShortTextField().
					SetDisplayName("Sheet Column Index").
					SetDescription("The index of the column where the new column should be added. Index is zero based.").
					SetRequired(true).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *AddColumnInWorkSheetOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[addColumnInWorkSheetOperationProps](ctx)
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
			SheetId:    convertToInt64(input.SheetID),
			Dimension:  "COLUMNS",
			StartIndex: convertToInt64(input.SheetColumnIndex),
			EndIndex:   convertToInt64(input.SheetColumnIndex) + 1,
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
		// Fields("id, name, mimeType, webViewLink, kind, createdTime").
		Do()
	return spreadsheet, err
}

func (c *AddColumnInWorkSheetOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *AddColumnInWorkSheetOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
