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

package jsonconverter

import (
	"github.com/wakflo/extensions/internal/integrations/jsonconverter/actions"
	"github.com/wakflo/go-sdk/integration"
)

type TextToJson struct{}

func (n *TextToJson) Auth() *integration.Auth {
	return &integration.Auth{
		Required: false,
	}
}

func (n *TextToJson) Triggers() []integration.Trigger {
	return []integration.Trigger{}
}

func (n *TextToJson) Actions() []integration.Action {
	return []integration.Action{
		actions.NewConvertToJsonAction(),
	}
}

func NewTextToJson() integration.Integration {
	return &TextToJson{}
}

var Integration = integration.Register(NewTextToJson())
