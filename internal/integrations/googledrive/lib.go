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

package googledrive

import (
	"github.com/wakflo/extensions/internal/integrations/googledrive/actions"
	"github.com/wakflo/extensions/internal/integrations/googledrive/shared"
	"github.com/wakflo/extensions/internal/integrations/googledrive/triggers"
	"github.com/wakflo/go-sdk/integration"
)

type GoogleDrive struct {
}

func (n *GoogleDrive) Auth() *integration.Auth {
	return &integration.Auth{
		Required: true,
		Schema:   *shared.SharedGoogleDriveAuth,
	}
}

func (n *GoogleDrive) Triggers() []integration.Trigger {
	return []integration.Trigger{
		triggers.NewNewFileTrigger(),
		triggers.NewNewFolderTrigger(),
	}
}

func (n *GoogleDrive) Actions() []integration.Action {
	return []integration.Action{
		actions.NewCreateFileAction(),
		actions.NewCreateFolderAction(),
	}
}

func NewGoogleDrive() integration.Integration {
	return &GoogleDrive{}
}

var Integration = integration.Register(NewGoogleDrive())
