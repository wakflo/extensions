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
	"github.com/wakflo/go-sdk/sdk"
)

type Smartsheet struct{}

func (n *Smartsheet) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SmartsheetSharedAuth,
	}
}

func (n *Smartsheet) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *Smartsheet) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewListSheetAction(),
	}
}

func NewSmartsheet() sdk.Integration {
	return &Smartsheet{}
}

var Integration = sdk.Register(NewSmartsheet())
