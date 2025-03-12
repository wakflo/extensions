package instagram

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/instagram/actions"
	"github.com/wakflo/extensions/internal/integrations/instagram/shared"
	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewInstagram(), Flow, ReadME)

type Instagram struct{}

func (n *Instagram) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: false,
		Schema:   *shared.InstagramSharedAuth,
	}
}

func (n *Instagram) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *Instagram) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewTagLocationAction(),

		actions.NewInviteCollaboratorAction(),

		actions.NewPostReelAction(),

		actions.NewPostMediaAction(),

		actions.NewPostSingleMediaAction(),
	}
}

func NewInstagram() sdk.Integration {
	return &Instagram{}
}
