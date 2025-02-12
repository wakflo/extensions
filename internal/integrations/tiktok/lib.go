package tiktok

import (
	"github.com/wakflo/extensions/internal/integrations/tiktok/actions"
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewTikTok())

type TikTok struct{}

func (n *TikTok) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: false,
	}
}

func (n *TikTok) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *TikTok) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewUploadMediaFromUrlAction(),

		actions.NewUploadMediaAction(),

		actions.NewPostMediaFromUrlAction(),

		actions.NewPostMediaAction(),
	}
}

func NewTikTok() sdk.Integration {
	return &TikTok{}
}
