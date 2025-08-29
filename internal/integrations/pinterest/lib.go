package pinterest

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/pinterest/actions"
	"github.com/wakflo/extensions/internal/integrations/pinterest/shared"
	"github.com/wakflo/extensions/internal/integrations/pinterest/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewPinterest())

type Pinterest struct{}

func (n *Pinterest) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Pinterest) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedPinterestAuth,
	}
}

func (n *Pinterest) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewPinCreatedTrigger(),
	}
}

func (n *Pinterest) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewGetPinAction(),
		actions.NewDeletePinAction(),
		// actions.NewSearchPinsAction(),
		actions.NewUpdatePinAction(),
		actions.NewGetPinAnalyticsAction(),
		// actions.NewCreatePinAction(),
	}
}

func NewPinterest() sdk.Integration {
	return &Pinterest{}
}
