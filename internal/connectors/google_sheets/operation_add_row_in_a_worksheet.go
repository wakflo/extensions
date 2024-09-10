package googlesheets

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type addRowInWorkSheetOperationProps struct {
	SpreadSheetID string          `json:"spreadSheetId"`
	SheetTitle    string          `json:"sheetTitle"`
	SheetRow      string          `json:"sheetRow"`
	Values        [][]interface{} `json:"values"`
}

type AddRowInWorkSheetOperation struct {
	options *sdk.OperationInfo
}

// func returnEmptyArray(input *addRowInWorkSheetOperationProps) [][]interface{} {
// 	if input.Values == nil || len(input.Values) == 0 {
// 		return [][]interface{}{}
// 	}
// 	return input.Values
// }

func NewAddRowInWorkSheetOperation() *AddRowInWorkSheetOperation {
	return &AddRowInWorkSheetOperation{
		options: &sdk.OperationInfo{
			Name:        "Add Row In Sheet",
			Description: "Add a row in a sheet",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"spreadsheetId": getSpreadsheetsInput("Spreadsheet", "spreadsheet ID", true),
				"sheetTitle":    getSheetTitleInput("Sheet", "select sheet", true),
				"sheetRow": autoform.NewShortTextField().
					SetDisplayName("Sheet Row").
					SetDescription("The row range of the sheet in the format of A1 notation where you want to append the new row.").
					SetRequired(true).
					Build(),
				"values": autoform.NewArrayField().
					SetDisplayName("Values").
					SetItems(
						autoform.NewShortTextField().
							SetDisplayName("Label").
							SetDescription("Label").
							SetRequired(true).
							Build(),
					).
					SetDescription("The values to be added in the row you're about to create.").
					SetRequired(false).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *AddRowInWorkSheetOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[addRowInWorkSheetOperationProps](ctx)
	sheetService, err := sheets.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	if input.SpreadSheetID == "" {
		return nil, errors.New("spreadsheet ID is required")
	}

	if input.SheetTitle == "" {
		return nil, errors.New("sheet Title is required")
	}

	if input.SheetRow == "" {
		return nil, errors.New("sheet row range is required")
	}

	// if len(input.Values) <= 0 {
	// 	return nil, errors.New("values are required")
	// }

	// this variable concatenates the sheet ID and the sheet row range since we want to append a row in a particular range of the worksheet
	// range_ := input.SheetID + input.SheetRow
	range_ := fmt.Sprintf("%s!%s", input.SheetTitle, input.SheetRow)

	// sample array data to add as new row
	sampleData := [][]interface{}{{"apple", "banana", "orange", "pineapple", "mango"}}

	spreadsheet, err := sheetService.Spreadsheets.Values.Append(input.SpreadSheetID, range_, &sheets.ValueRange{
		Range:          range_,
		MajorDimension: "ROWS",
		Values:         sampleData,
	}).
		ValueInputOption("RAW").
		InsertDataOption("INSERT_ROWS").
		Do()
	return spreadsheet, err
}

func (c *AddRowInWorkSheetOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *AddRowInWorkSheetOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
