package gemini

import (
	"github.com/wakflo/extensions/internal/integrations/gemini/actions"
	"github.com/wakflo/go-sdk/autoform"
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewGemini())
var sharedAuth = autoform.NewAuthSecretField().SetDisplayName("Gemini API key").SetDescription("Your Gemini api key").Build()

type Gemini struct{}

func (n *Gemini) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *sharedAuth,
	}
}

func (n *Gemini) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *Gemini) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewChatGeminiAction(),
	}
}

func NewGemini() sdk.Integration {
	return &Gemini{}
}
