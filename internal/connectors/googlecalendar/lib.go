package googlecalendar

import (
	sdk "github.com/wakflo/go-sdk/connector"
)

func NewConnector() (*sdk.ConnectorPlugin, error) {
	return sdk.CreateConnector(&sdk.CreateConnectorArgs{
		Name:        "Google Calendar",
		Description: "some google calendar connector",
		Logo:        "logos:google-calendar",
		Version:     "0.0.1",
		Group:       sdk.ConnectorGroupApps,
		Authors:     []string{"Wakflo <integrations@wakflo.com>"},
		Triggers: []sdk.ITrigger{
			NewTriggerNewEventCreated(),
		},
		Operations: []sdk.IOperation{
			NewCreateEventOperation(),
			NewUpdateEventOperation(),
		},
	})
}
