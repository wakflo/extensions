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

type updateRowInWorkSheetOperationProps struct {
	SpreadSheetID string          `json:"spreadSheetId"`
	SheetTitle    string          `json:"sheetTitle"`
	SheetRow      string          `json:"sheetRow"`
	Values        [][]interface{} `json:"values"`
}

type UpdateRowInWorkSheetOperation struct {
	options *sdk.OperationInfo
}

func NewUpdateRowInWorkSheetOperation() *UpdateRowInWorkSheetOperation {
	return &UpdateRowInWorkSheetOperation{
		options: &sdk.OperationInfo{
			Name:        "Update Row In Sheet",
			Description: "Update a row in a sheet",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"spreadSheetId": autoform.NewShortTextField().
					SetDisplayName("Spreadsheet ID").
					SetDescription("The ID of the spreadsheet.").
					SetRequired(true).
					Build(),
				"sheetTitle": autoform.NewShortTextField().
					SetDisplayName("Sheet Title").
					SetDescription("The Title of the sheet.").
					SetRequired(true).
					Build(),
				"sheetRow": autoform.NewShortTextField().
					SetDisplayName("Sheet Row").
					SetDescription("The row range of the sheet in the format of A1 notation.").
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
					SetDescription("The values to be updated in the row.").
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

func (c *UpdateRowInWorkSheetOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[updateRowInWorkSheetOperationProps](ctx)
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

	// this variable concatenates the sheet ID and the sheet row range since we want to update a particular row in a sheet
	// range_ := input.SheetID + input.SheetRow
	range_ := fmt.Sprintf("%s!%s", input.SheetTitle, input.SheetRow)

	// sample array data to add as new row
	sampleData := [][]interface{}{{"john", "mary", "sophia", "james", "victor"}}

	spreadsheet, err := sheetService.Spreadsheets.Values.Update(input.SpreadSheetID, range_, &sheets.ValueRange{
		Range:          range_,
		MajorDimension: "ROWS",
		Values:         sampleData,
	}).
		ValueInputOption("RAW").
		// Fields("id, name, mimeType, webViewLink, kind, createdTime").
		Do()
	return spreadsheet, err
}

func (c *UpdateRowInWorkSheetOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *UpdateRowInWorkSheetOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
