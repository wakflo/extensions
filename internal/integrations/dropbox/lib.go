package dropbox

import (
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewDropbox())

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
