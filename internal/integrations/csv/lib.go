package csv

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/csv/actions"
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewCSV(), Flow, ReadME)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

type CSV struct{}

func (n *CSV) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: false,
	}
}

func (n *CSV) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *CSV) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewRowCountAction(),
	}
}

func NewCSV() sdk.Integration {
	return &CSV{}
}
