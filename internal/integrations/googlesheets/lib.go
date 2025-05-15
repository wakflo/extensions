package googlesheets

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/googlesheets/actions"
	"github.com/wakflo/extensions/internal/integrations/googlesheets/shared"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewGoogleSheets())

type GoogleSheets struct{}

func (n *GoogleSheets) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *GoogleSheets) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedGoogleSheetsAuth,
	}
}

func (n *GoogleSheets) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *GoogleSheets) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewUpdateRowInWorksheetAction(),

		actions.NewReadRowInWorksheetAction(),

		actions.NewFindWorksheetAction(),

		actions.NewGetWorksheetByIdAction(),

		actions.NewCopyWorksheetAction(),

		actions.NewCreateSpreadsheetAction(),

		actions.NewAddRowInWorksheetAction(),

		actions.NewAddColumnInWorksheetAction(),
	}
}

func NewGoogleSheets() sdk.Integration {
	return &GoogleSheets{}
}
