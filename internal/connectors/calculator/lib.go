package calculator

import (
	sdk "github.com/wakflo/go-sdk/connector"
)

func NewConnector() (*sdk.ConnectorPlugin, error) {
	return sdk.CreateConnector(&sdk.CreateConnectorArgs{
		Name:        "Calculator",
		Description: "Make simple math operations",
		Logo:        "ion:calculator",
		Version:     "0.0.1",
		Category:    sdk.Core,
		Authors:     []string{"Wakflo <integrations@wakflo.com>"},
		Triggers:    []sdk.ITrigger{},
		Operations: []sdk.IOperation{
			NewAddition(),
			NewSubtraction(),
			NewDivision(),
			NewMultiplication(),
			NewModulo(),
		},
	})
}
