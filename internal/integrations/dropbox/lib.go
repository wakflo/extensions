package dropbox

import (
	_ "embed"

	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewDropbox(), Flow, ReadME)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

type Dropbox struct{}

func (n *Dropbox) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: false,
	}
}

func (n *Dropbox) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *Dropbox) Actions() []sdk.Action {
	return []sdk.Action{}
}

func NewDropbox() sdk.Integration {
	return &Dropbox{}
}
