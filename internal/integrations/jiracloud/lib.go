package jiracloud

import (
	_ "embed"

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
	return []sdk.Trigger{}
}

func (n *JiraCloud) Actions() []sdk.Action {
	return []sdk.Action{}
}

func NewJiraCloud() sdk.Integration {
	return &JiraCloud{}
}
