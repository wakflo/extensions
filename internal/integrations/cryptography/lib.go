package cryptography

import (
	"github.com/wakflo/extensions/internal/integrations/cryptography/actions"
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewCryptography())

type Cryptography struct{}

func (n *Cryptography) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: false,
	}
}

func (n *Cryptography) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *Cryptography) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewHashTextAction(),

		actions.NewGenerateTextAction(),
	}
}

func NewCryptography() sdk.Integration {
	return &Cryptography{}
}
