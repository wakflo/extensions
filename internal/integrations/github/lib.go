package github

import (
	"github.com/wakflo/extensions/internal/integrations/github/actions"
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewGithub())

type Github struct{}

func (n *Github) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: false,
	}
}

func (n *Github) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *Github) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewUnlockIssueAction(),

		actions.NewLockIssueAction(),

		actions.NewGetIssueAction(),

		actions.NewCreateIssueAction(),

		actions.NewCreateIssueCommentAction(),
	}
}

func NewGithub() sdk.Integration {
	return &Github{}
}
