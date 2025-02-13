package googlesheets

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/googlesheets/actions"
	"github.com/wakflo/extensions/internal/integrations/googlesheets/shared"
	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewGoogleSheets(), Flow, ReadME)

type GoogleSheets struct{}

func (n *GoogleSheets) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedAuth,
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
