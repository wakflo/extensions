package jiracloud

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/jiracloud/actions"
	"github.com/wakflo/extensions/internal/integrations/jiracloud/triggers"
	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewJiraCloud(), Flow, ReadME)

type JiraCloud struct{}

func (n *JiraCloud) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: false,
	}
}

func (n *JiraCloud) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewIssueCreatedTrigger(),
		triggers.NewIssueUpdatedTrigger(),
	}
}

func (n *JiraCloud) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewAddCommentAction(),
		actions.NewCreateIssueAction(),
		actions.NewGetIssueAction(),
		actions.NewListIssuesAction(),
		actions.NewUpdateIssueAction(),
		actions.NewTransitionIssueAction(),
	}
}

func NewJiraCloud() sdk.Integration {
	return &JiraCloud{}
}
