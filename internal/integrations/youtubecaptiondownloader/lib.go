package youtubecaptiondownloader

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/youtubecaptiondownloader/actions"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewYoutubeCaptionDownloader())

type YoutubeCaptionDownloader struct{}

func (n *YoutubeCaptionDownloader) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *YoutubeCaptionDownloader) Auth() *core.AuthMetadata {
	return nil
}

func (n *YoutubeCaptionDownloader) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *YoutubeCaptionDownloader) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewGetVideoCaptionAction(),
	}
}

func NewYoutubeCaptionDownloader() sdk.Integration {
	return &YoutubeCaptionDownloader{}
}
