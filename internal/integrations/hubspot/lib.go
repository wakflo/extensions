package hubspot

import (
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewHubspot())

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
