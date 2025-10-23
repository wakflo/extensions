package socialkit

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/socialkit/actions"
	"github.com/wakflo/extensions/internal/integrations/socialkit/shared"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewSocialKit())

type SocialKit struct{}

func (n *SocialKit) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (s *SocialKit) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedAuth,
	}
}

func (s *SocialKit) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (s *SocialKit) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewGetYouTubeTranscriptAction(),
	}
}

func NewSocialKit() sdk.Integration {
	return &SocialKit{}
}
