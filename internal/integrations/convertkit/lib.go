package convertkit

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/convertkit/actions"
	"github.com/wakflo/extensions/internal/integrations/convertkit/shared"
	"github.com/wakflo/extensions/internal/integrations/convertkit/triggers"
	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewConvertKit(), Flow, ReadME)

type ConvertKit struct{}

func (n *ConvertKit) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedAuth,
	}
}

func (n *ConvertKit) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewSubscriberCreatedTrigger(),
		triggers.NewTagCreatedTrigger(),
	}
}

func (n *ConvertKit) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewListSubscribersAction(),
		actions.NewGetSubscriberAction(),
		// actions.NewCreateSubscriberAction(),
		actions.NewCreateTagAction(),
		actions.NewTagSubscriberAction(),
		actions.NewListTagsAction(),
	}
}

func NewConvertKit() sdk.Integration {
	return &ConvertKit{}
}
