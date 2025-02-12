package linear

import (
	"github.com/wakflo/extensions/internal/integrations/linear/actions"
	"github.com/wakflo/extensions/internal/integrations/linear/shared"
	"github.com/wakflo/extensions/internal/integrations/linear/triggers"
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewLinear())

type Linear struct{}

func (n *Linear) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedAuth,
	}
}

func (n *Linear) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewIssueUpdatedTrigger(),

		triggers.NewIssueCreatedTrigger(),
	}
}

func (n *Linear) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewUpdateIssueAction(),

		actions.NewFindIssuesAction(),

		actions.NewCreateIssueAction(),
	}
}

func NewLinear() sdk.Integration {
	return &Linear{}
}
