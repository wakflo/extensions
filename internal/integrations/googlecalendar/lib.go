package googlecalendar

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/googlecalendar/actions"
	"github.com/wakflo/extensions/internal/integrations/googlecalendar/shared"
	"github.com/wakflo/extensions/internal/integrations/googlecalendar/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewGoogleCalendar())

type GoogleCalendar struct{}

func (n *GoogleCalendar) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *GoogleCalendar) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedGoogleCalendarAuth,
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
