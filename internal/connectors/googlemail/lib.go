package googlemail

import (
	sdk "github.com/wakflo/go-sdk/connector"
)

func NewConnector() (*sdk.ConnectorPlugin, error) {
	return sdk.CreateConnector(&sdk.CreateConnectorArgs{
		Name:        "Google Mail",
		Description: " Connect to Google Mail to read, send and manage emails",
		Logo:        "logos:google-gmail",
		Version:     "0.0.1",
		Category:    sdk.Apps,
		Authors:     []string{"Wakflo <integrations@wakflo.com>"},
		Triggers: []sdk.ITrigger{
			NewTriggerNewEmail(),
		},
		Operations: []sdk.IOperation{
			NewGetMailByIDOperation(),
			NewSendMailOperation(),
			NewGetThreadOperation(),
			NewSendTemplateMailOperation(),
			NewListMailsOperation(),
		},
	})
}
