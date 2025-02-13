package cryptography

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/cryptography/actions"
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewCryptography(), Flow, ReadME)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

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
