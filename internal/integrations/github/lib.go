package github

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/github/actions"
	"github.com/wakflo/extensions/internal/integrations/github/shared"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewGithub())

type Github struct{}

func (n *Github) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Github) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedGithubAuth,
	}
}

func (n *Github) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *Github) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewCreateIssueAction(),
		actions.NewLockIssueAction(),
		actions.NewUnlockIssueAction(),
		actions.NewGetIssueAction(),
		actions.NewCreateIssueCommentAction(),
	}
}

func NewGithub() sdk.Integration {
	return &Github{}
}
