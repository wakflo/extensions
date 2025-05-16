package airtable

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/airtable/actions"
	"github.com/wakflo/extensions/internal/integrations/airtable/shared"
	"github.com/wakflo/extensions/internal/integrations/airtable/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewAirtable())

type Airtable struct{}

func (n *Airtable) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Airtable) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.AirtableSharedAuth,
	}
}

func (n *Airtable) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewNewRecordTrigger(),
	}
}

func (n *Airtable) Actions() []sdk.Action {
	return []sdk.Action{
		// actions.NewUpdateRecordAction(),
		actions.NewFindRecordAction(),
		actions.NewDeleteRecordAction(),
	}
}

func NewAirtable() sdk.Integration {
	return &Airtable{}
}
