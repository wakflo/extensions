package dropbox

import (
	_ "embed"

	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewDropbox())

type Dropbox struct{}

func (n *Dropbox) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Dropbox) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: false,
	}
}

func (n *Dropbox) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *Dropbox) Actions() []sdk.Action {
	return []sdk.Action{}
}

func NewDropbox() sdk.Integration {
	return &Dropbox{}
}
