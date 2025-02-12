package instagram

import (
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewInstagram())

type Instagram struct{}

func (n *Instagram) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: false,
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
