package jiracloudsoftware

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/jiracloudsoftware/actions"
	"github.com/wakflo/extensions/internal/integrations/jiracloudsoftware/shared"
	"github.com/wakflo/extensions/internal/integrations/jiracloudsoftware/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewJiraCloudSoftware())

type JiraCloudSoftware struct{}

func (n *JiraCloudSoftware) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *JiraCloudSoftware) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.JiraSharedAuth,
	}
}

func (n *JiraCloudSoftware) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewIssueCreatedTrigger(),
		triggers.NewIssueUpdatedTrigger(),
	}
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
