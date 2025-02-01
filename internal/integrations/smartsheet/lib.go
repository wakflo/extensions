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

package smartsheet

import (
	"github.com/wakflo/extensions/internal/integrations/smartsheet/actions"
	"github.com/wakflo/extensions/internal/integrations/smartsheet/shared"
	"github.com/wakflo/go-sdk/integration"
)

type Smartsheet struct{}

func (n *Smartsheet) Auth() *integration.Auth {
	return &integration.Auth{
		Required: true,
		Schema:   *shared.SmartsheetSharedAuth,
	}
}

func (n *Smartsheet) Triggers() []integration.Trigger {
	return []integration.Trigger{}
}

func (n *Smartsheet) Actions() []integration.Action {
	return []integration.Action{
		actions.NewListSheetAction(),
	}
}

func NewSmartsheet() integration.Integration {
	return &Smartsheet{}
}

var Integration = integration.Register(NewSmartsheet())
