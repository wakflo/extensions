package gemini

import (
	_ "embed"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/gemini/actions"
	"github.com/wakflo/go-sdk/v2"

	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

// var sharedAuth = smartform.NewAuthForm().SetDisplayName("Gemini API key").SetDescription("Your Gemini api key").Build()

var (
	form = smartform.NewAuthForm("gemini-auth", "Gemini API Authentication", smartform.AuthStrategyCustom)

	_ = form.TextField("key", "Gemini API Key").
		Required(true).
		HelpText("Your Gemini api key")

	GeminiSharedAuth = form.Build()
)

var Integration = sdk.Register(NewGemini())

type Gemini struct{}

func (n *Gemini) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Gemini) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   GeminiSharedAuth,
	}
}

func (n *Gemini) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *Gemini) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewChatGeminiAction(),
		actions.NewGenerateEmbeddingAction(),
		actions.NewTranslateTextAction(),
		actions.NewSummarizeTextAction(),
		actions.NewAnalyzeImageAction(),
		actions.NewFunctionCallingAction(),
	}
}

func NewGemini() sdk.Integration {
	return &Gemini{}
}
