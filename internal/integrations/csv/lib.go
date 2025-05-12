package csv

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/csv/actions"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewCSV())

type CSV struct{}

func (n *CSV) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *CSV) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
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
