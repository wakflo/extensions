package typeform

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/typeform/actions"
	"github.com/wakflo/extensions/internal/integrations/typeform/shared"
	"github.com/wakflo/extensions/internal/integrations/typeform/triggers"
	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewTypeform(), Flow, ReadME)

type Typeform struct{}

func (t *Typeform) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedAuth,
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
