package linear

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/linear/actions"
	"github.com/wakflo/extensions/internal/integrations/linear/shared"
	"github.com/wakflo/extensions/internal/integrations/linear/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

var (
	//go:embed README.md
	ReadME string

	//go:embed flo.toml
	Flow string

	Integration = sdk.Register(NewLinear())
)

type Linear struct{}

func (n *Linear) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Linear) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedAuth,
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
