package freshworkscrm

import (
	_ "embed"

	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewFreshWorksCRM(), Flow, ReadME)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

type FreshWorksCRM struct{}

func (n *FreshWorksCRM) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: false,
	}
}

func (n *FreshWorksCRM) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *FreshWorksCRM) Actions() []sdk.Action {
	return []sdk.Action{}
}

func NewFreshWorksCRM() sdk.Integration {
	return &FreshWorksCRM{}
}
