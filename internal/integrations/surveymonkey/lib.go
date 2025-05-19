package surveyMonkey

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/surveymonkey/actions"
	"github.com/wakflo/extensions/internal/integrations/surveymonkey/shared"
	"github.com/wakflo/extensions/internal/integrations/surveymonkey/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewSurveyMonkey())

type SurveyMonkey struct{}

func (t *SurveyMonkey) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (t *SurveyMonkey) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SurveyMonkeySharedAuth,
	}
}

func (t *SurveyMonkey) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewNewResponseTrigger(),
	}
}

func (t *SurveyMonkey) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewCreateContactListAction(),
	}
}

func NewSurveyMonkey() sdk.Integration {
	return &SurveyMonkey{}
}
