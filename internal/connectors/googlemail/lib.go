package googlemail

import (
	_ "embed"

	sdk "github.com/wakflo/go-sdk/connector"
)

//go:embed docs.mdx
var doc string

func NewConnector() (*sdk.ConnectorPlugin, error) {
	return sdk.CreateConnector(&sdk.CreateConnectorArgs{
		Name:          "Google Mail",
		Description:   "Connect to Google Mail to read, send and manage emails",
		Logo:          "logos:google-gmail",
		Version:       "0.0.1",
		Group:         sdk.ConnectorGroupApps,
		Authors:       []string{"Wakflo <integrations@wakflo.com>"},
		Documentation: doc,
		Triggers: []sdk.ITrigger{
			NewTriggerNewEmail(),
		},
		Operations: []sdk.IOperation{
			NewGetMailByIDOperation(),
			NewSendMailOperation(),
			NewGetThreadOperation(),
			NewListMailsOperation(),
		},
	})
}
