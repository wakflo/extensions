package zoom

import (
	"github.com/wakflo/extensions/internal/integrations/zoom/actions"
	"github.com/wakflo/extensions/internal/integrations/zoom/shared"
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewZoom())

type Zoom struct{}

func (n *Zoom) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.ZoomSharedAuth,
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
