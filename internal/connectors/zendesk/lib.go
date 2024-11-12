package zendesk

import (
	sdk "github.com/wakflo/go-sdk/connector"
)

func NewConnector() (*sdk.ConnectorPlugin, error) {
	return sdk.CreateConnector(&sdk.CreateConnectorArgs{
		Name:        "Zendesk",
		Description: "Customer service software and support ticket system",
		Logo:        "logos:zendesk-icon",
		Version:     "0.0.1",
		Group:       sdk.ConnectorGroupApps,
		Authors:     []string{"Wakflo <integrations@wakflo.com>"},
		Triggers: []sdk.ITrigger{
			NewTicketCreated(),
		},
		Operations: []sdk.IOperation{
			NewGetGroupsOperation(),
			NewGetTicketsOperation(),
		},
	})
}
