package typeform

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/typeform/actions"
	"github.com/wakflo/extensions/internal/integrations/typeform/shared"
	"github.com/wakflo/extensions/internal/integrations/typeform/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewTypeform())

type Typeform struct{}

func (t *Typeform) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (t *Typeform) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedTypeformAuth,
	}
}

func (t *Typeform) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewNewResponseTrigger(),
	}
}

func (t *Typeform) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewRetrieveFormResponseAction(),
	}
}

func NewTypeform() sdk.Integration {
	return &Typeform{}
}
