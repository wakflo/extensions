package zohocrm

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/zohocrm/actions"
	"github.com/wakflo/extensions/internal/integrations/zohocrm/shared"
	"github.com/wakflo/extensions/internal/integrations/zohocrm/triggers"

	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewZohoCrm())

type ZohoCrm struct{}

func (n *ZohoCrm) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *ZohoCrm) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedAuth,
	}
}

func (n *ZohoCrm) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewRecordUpdatedTrigger(),
		triggers.NewRecordCreatedTrigger(),
	}
}

func (n *ZohoCrm) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewCreateRecordAction(),
		actions.NewUpdateRecordAction(),
		actions.NewGetRecordAction(),
		actions.NewSearchRecordsAction(),
		actions.NewDeleteRecordAction(),
		actions.NewListRecordsAction(),
	}
}

func NewZohoCrm() sdk.Integration {
	return &ZohoCrm{}
}
