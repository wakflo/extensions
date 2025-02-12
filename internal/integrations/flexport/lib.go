package flexport

import (
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewFlexport())

type Flexport struct{}

func (n *Flexport) Auth() *sdk.Auth {
	return &sdk.Auth{
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
