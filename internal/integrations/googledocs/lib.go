package googledocs

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/googledocs/actions"
	"github.com/wakflo/extensions/internal/integrations/googledocs/shared"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewGoogleDocs())

type GoogleDocs struct{}

func (n *GoogleDocs) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *GoogleDocs) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedGoogleDocsAuth,
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
