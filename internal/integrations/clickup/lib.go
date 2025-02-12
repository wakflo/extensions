package clickup

import (
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewClickup())

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
