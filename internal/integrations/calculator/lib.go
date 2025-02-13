package calculator

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/calculator/actions"
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewCalculator(), Flow, ReadME)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

type Calculator struct{}

func (n *Calculator) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: false,
	}
}

func (n *Calculator) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *Calculator) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewAddAction(),
		actions.NewSubtractAction(),
		actions.NewMultiplyAction(),
		actions.NewDivideAction(),
		actions.NewModuloAction(),
	}
}

func NewCalculator() sdk.Integration {
	return &Calculator{}
}
