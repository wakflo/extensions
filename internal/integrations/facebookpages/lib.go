// Copyright 2022-present Wakflo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package facebookpages

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/facebookpages/actions"
	"github.com/wakflo/extensions/internal/integrations/facebookpages/shared"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewFacebookPages())

type FacebookPages struct{}

func (n *FacebookPages) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *FacebookPages) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.FacebookPagesSharedAuth,
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
