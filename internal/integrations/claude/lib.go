package claude

import (
	_ "embed"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/claude/actions"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var (
	form = smartform.NewAuthForm("claude-auth", "Claude API Authentication", smartform.AuthStrategyCustom)

	_ = form.TextField("apiKey", "API Key").
		Required(true).
		Placeholder("sk-ant-api03-...").
		HelpText("You can get your API key from https://console.anthropic.com/settings/keys")

	ClaudeSharedAuth = form.Build()
)

var Integration = sdk.Register(NewClaude())

type Claude struct{}

func (n *Claude) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (c *Claude) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   ClaudeSharedAuth,
	}
}

func (c *Claude) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (c *Claude) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewChatClaudeAction(),
		actions.NewChatWithImagesAction(),
		actions.NewSummarizeTextAction(),
		actions.NewTranslateTextAction(),
		// actions.NewExtractStructuredDataAction(),

		// actions.NewGenerateCodeAction(),

		actions.NewAnalyzeTextAction(),
		actions.NewCompareTextsAction(),
	}
}

func NewClaude() sdk.Integration {
	return &Claude{}
}
