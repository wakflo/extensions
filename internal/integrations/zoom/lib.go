package zoom

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/zoom/actions"
	"github.com/wakflo/extensions/internal/integrations/zoom/shared"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewZoom())

type Zoom struct{}

func (n *Zoom) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Zoom) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.ZoomSharedAuth,
	}
}

func (n *Zoom) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *Zoom) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewCreateMeetingRegistrationAction(),
		actions.NewCreateMeetingAction(),
	}
}

func NewZoom() sdk.Integration {
	return &Zoom{}
}
