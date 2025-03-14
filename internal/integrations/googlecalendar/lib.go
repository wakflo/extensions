package googlecalendar

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/googlecalendar/actions"
	"github.com/wakflo/extensions/internal/integrations/googlecalendar/shared"
	"github.com/wakflo/extensions/internal/integrations/googlecalendar/triggers"
	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewGoogleCalendar(), Flow, ReadME)

type GoogleCalendar struct{}

func (n *GoogleCalendar) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedAuth,
	}
}

func (n *GoogleCalendar) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewEventCreatedTrigger(),
	}
}

func (n *GoogleCalendar) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewUpdateEventAction(),

		actions.NewCreateEventAction(),
	}
}

func NewGoogleCalendar() sdk.Integration {
	return &GoogleCalendar{}
}
