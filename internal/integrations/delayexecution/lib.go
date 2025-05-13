package delayexecution

import (
	_ "embed"

	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewDelayExecution())

type DelayExecution struct{}

func (n *DelayExecution) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *DelayExecution) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
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
