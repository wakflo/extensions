package cryptography

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/cryptography/actions"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewCryptography())

type Cryptography struct{}

func (n *Cryptography) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Cryptography) Auth() *core.AuthMetadata {
	return nil
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
