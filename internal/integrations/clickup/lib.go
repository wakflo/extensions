package clickup

import (
	_ "embed"

	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewClickup(), Flow, ReadME)

type Clickup struct{}

func (n *Clickup) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: false,
	}
}

func (n *Clickup) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *Clickup) Actions() []sdk.Action {
	return []sdk.Action{}
}

func NewClickup() sdk.Integration {
	return &Clickup{}
}

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string
