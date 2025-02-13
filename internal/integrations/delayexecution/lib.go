package delayexecution

import (
	_ "embed"

	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewDelayExecution(), Flow, ReadME)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

type DelayExecution struct{}

func (n *DelayExecution) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: false,
	}
}

func (n *DelayExecution) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *DelayExecution) Actions() []sdk.Action {
	return []sdk.Action{}
}

func NewDelayExecution() sdk.Integration {
	return &DelayExecution{}
}
