package calculator

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/calculator/actions"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewCalculator())

type Calculator struct{}

func (n *Calculator) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Calculator) Auth() *core.AuthMetadata {
	return nil
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
