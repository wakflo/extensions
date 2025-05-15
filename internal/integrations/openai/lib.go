package openai

import (
	_ "embed"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/openai/actions"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var sharedAuth = smartform.NewAuthForm("openai", "OpenAI", smartform.AuthStrategyCustom)

var Integration = sdk.Register(NewOpenAI())

type OpenAI struct{}

func (n *OpenAI) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *OpenAI) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   actions.OpenAISharedAuth,
	}
}

func (n *OpenAI) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *OpenAI) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewChatOpenAIAction(),
	}
}

func NewOpenAI() sdk.Integration {
	return &OpenAI{}
}
