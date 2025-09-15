package whatsapp

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/whatsapp/actions"
	"github.com/wakflo/extensions/internal/integrations/whatsapp/shared"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewWhatsApp())

type WhatsApp struct{}

func (w *WhatsApp) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (w *WhatsApp) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedAuth,
	}
}

func (w *WhatsApp) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (w *WhatsApp) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewSendMessageAction(),
		actions.NewSendMediaAction(),
	}
}

func NewWhatsApp() sdk.Integration {
	return &WhatsApp{}
}
