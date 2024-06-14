package googlesheets

import (
	sdk "github.com/wakflo/go-sdk/connector"
)

func NewConnector() (*sdk.ConnectorPlugin, error) {
	return sdk.CreateConnector(&sdk.CreateConnectorArgs{
		Name:        "Google Sheets",
		Description: "some google sheets connector",
		Logo:        "hugeicons:google-sheet",
		Version:     "0.0.1",
		Category:    sdk.Apps,
		Authors:     []string{"Wakflo <integrations@wakflo.com>"},
		Triggers:    []sdk.ITrigger{},
		Operations: []sdk.IOperation{
			NewCreateSheetOperation(),
			NewUpdateRowInWorkSheetOperation(),
			NewAddRowInWorkSheetOperation(),
			NewCopyWorkSheetOperation(),
			NewAddColumnInWorkSheetOperation(),
			NewFindWorkSheetByIDOperation(),
			NewFindWorkSheetByTitleOperation(),
		},
	})
}
