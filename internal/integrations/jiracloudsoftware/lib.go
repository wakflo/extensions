package jiracloudsoftware

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/jiracloudsoftware/actions"
	"github.com/wakflo/extensions/internal/integrations/jiracloudsoftware/shared"
	// "github.com/wakflo/extensions/internal/integrations/jiracloudsoftware/triggers"
	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewJiraCloudSoftware(), Flow, ReadME)

type JiraCloudSoftware struct{}

func (n *JiraCloudSoftware) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedAuth,
	}
}

func (n *JiraCloudSoftware) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *JiraCloudSoftware) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewAddCommentAction(),
		actions.NewCreateIssueAction(),
		actions.NewGetIssueAction(),
		actions.NewListIssuesAction(),
		actions.NewUpdateIssueAction(),
		actions.NewTransitionIssueAction(),
	}
}

func NewJiraCloudSoftware() sdk.Integration {
	return &JiraCloudSoftware{}
}
