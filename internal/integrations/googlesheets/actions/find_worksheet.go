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
			"message": "Hello World!",
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

	spreadsheet, err := sheetService.Spreadsheets.Get(input.SpreadSheetID).
		Do()
	if err != nil {
		return nil, err
	}

	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.Title == input.SheetTitle {
			return sheet, nil
		}

		if sheet.Properties.Title != input.SheetTitle {
			const notFound string = "not sheet found"
			return notFound, nil
		}
	}
	return spreadsheet, err
}

func NewFindWorksheetAction() sdk.Action {
	return &FindWorksheetAction{}
}
