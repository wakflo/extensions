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
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/jsonconverter/actions"
	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

type TextToJson struct{}

func (n *TextToJson) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: false,
	}
}

func (n *TextToJson) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *TextToJson) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewConvertToJsonAction(),
	}
}

func NewTextToJson() sdk.Integration {
	return &TextToJson{}
}

var Integration = sdk.Register(NewTextToJson(), Flow, ReadME)
