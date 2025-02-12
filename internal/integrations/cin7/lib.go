package cin7

import (
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewCin7())

type Cin7 struct{}

func (n *Cin7) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: false,
	}
}

func (n *Cin7) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *Cin7) Actions() []sdk.Action {
	return []sdk.Action{}
}

func NewCin7() sdk.Integration {
	return &Cin7{}
}
