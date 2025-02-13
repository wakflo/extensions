package googledocs

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/googledocs/actions"
	"github.com/wakflo/extensions/internal/integrations/googledocs/shared"
	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewGoogleDocs(), Flow, ReadME)

type GoogleDocs struct{}

func (n *GoogleDocs) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedAuth,
	}
}

func (n *GoogleDocs) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *GoogleDocs) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewReadDocumentAction(),

		actions.NewFindDocumentAction(),

		actions.NewCreateDocumentAction(),

		actions.NewAppendTextToDocumentAction(),
	}
}

func NewGoogleDocs() sdk.Integration {
	return &GoogleDocs{}
}
