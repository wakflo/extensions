package hubspot

import (
	_ "embed"

	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewHubspot(), Flow, ReadME)

type Hubspot struct{}

func (n *Hubspot) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: false,
	}
}

func (n *Hubspot) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *Hubspot) Actions() []sdk.Action {
	return []sdk.Action{}
}

func NewHubspot() sdk.Integration {
	return &Hubspot{}
}
