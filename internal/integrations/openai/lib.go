package openai

import (
	"github.com/wakflo/extensions/internal/integrations/openai/actions"
	"github.com/wakflo/go-sdk/autoform"
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewOpenAI())

var sharedAuth = autoform.NewAuthSecretField().SetDisplayName("OpenAI API key").SetDescription("Your OpenAI api key").Build()

type OpenAI struct{}

func (n *OpenAI) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *sharedAuth,
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
