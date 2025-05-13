package flexport

import (
	_ "embed"

	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewFlexport())

type Flexport struct{}

func (n *Flexport) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Flexport) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: false,
	}
}

func (n *Flexport) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *Flexport) Actions() []sdk.Action {
	return []sdk.Action{}
}

func NewFlexport() sdk.Integration {
	return &Flexport{}
}
