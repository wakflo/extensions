package airtable

import (
	"github.com/wakflo/extensions/internal/integrations/airtable/actions"
	"github.com/wakflo/extensions/internal/integrations/airtable/shared"
	"github.com/wakflo/extensions/internal/integrations/airtable/triggers"
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewAirtable())

type Airtable struct{}

func (n *Airtable) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedAuth,
	}
}

func (n *Airtable) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewNewRecordTrigger(),
	}
}

func (n *Airtable) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewUpdateRecordAction(),
		actions.NewFindRecordAction(),
		actions.NewDeleteRecordAction(),
	}
}

func NewAirtable() sdk.Integration {
	return &Airtable{}
}
