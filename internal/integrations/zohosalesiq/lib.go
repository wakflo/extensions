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

package zohosalesiq

import (
	"github.com/wakflo/extensions/internal/integrations/zohosalesiq/actions"
	"github.com/wakflo/extensions/internal/integrations/zohosalesiq/shared"
	"github.com/wakflo/go-sdk/sdk"
)

type ZohoSalesIq struct{}

func (n *ZohoSalesIq) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.ZohoSalesSharedAuth,
	}
}

func (n *ZohoSalesIq) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *ZohoSalesIq) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewListChatsAction(),
		actions.NewGetVisitorsDetailsAction(),
	}
}

func NewZohoSalesIq() sdk.Integration {
	return &ZohoSalesIq{}
}

var Integration = sdk.Register(NewZohoSalesIq())
