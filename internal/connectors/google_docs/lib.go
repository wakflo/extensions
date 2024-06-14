package googledocs

import (
	sdk "github.com/wakflo/go-sdk/connector"
)

func NewConnector() (*sdk.ConnectorPlugin, error) {
	return sdk.CreateConnector(&sdk.CreateConnectorArgs{
		Name:        "Google Docs",
		Description: "some google docs connector",
		Logo:        "icon-park:google",
		Version:     "0.0.1",
		Category:    sdk.Apps,
		Authors:     []string{"Wakflo <integrations@wakflo.com>"},
		Triggers:    []sdk.ITrigger{},
		Operations: []sdk.IOperation{
			NewFindDocumentOperation(),
			NewCreateDocumentOperation(),
			NewReadDocumentOperation(),
			NewAppendTextToDocumentOperation(),
		},
	})
}
