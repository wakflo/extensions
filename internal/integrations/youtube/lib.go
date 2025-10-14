package youtube

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/youtube/actions"
	"github.com/wakflo/extensions/internal/integrations/youtube/shared"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewYoutube())

type Youtube struct{}

func (n *Youtube) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Youtube) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedYoutubeAuuth,
	}
}

func (n *Youtube) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *Youtube) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewListVideosAction(),
		actions.NewGetVideoAction(),
		actions.NewUpdateVideoAction(),
		actions.NewUploadVideoAction(),
		actions.NewDownloadCaptionAction(),
	}
}

func NewYoutube() sdk.Integration {
	return &Youtube{}
}
