package facebookpages

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/facebookpages/actions"
	"github.com/wakflo/extensions/internal/integrations/facebookpages/shared"
	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewFacebookPages(), Flow, ReadME)

type FacebookPages struct{}

func (n *FacebookPages) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.FacebookPagesSharedAuth,
	}
}

func (n *FacebookPages) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *FacebookPages) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewGetPostsAction(),
		actions.NewDeletePostAction(),
		actions.NewPublishPostAction(),
		actions.NewUpdatePostAction(),
		actions.NewCreatePhotoPostAction(),
		actions.NewCreateVideoPostAction(),
	}
}

func NewFacebookPages() sdk.Integration {
	return &FacebookPages{}
}
