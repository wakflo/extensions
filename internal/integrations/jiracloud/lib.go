package jiracloud

import (
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewJiraCloud())

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
