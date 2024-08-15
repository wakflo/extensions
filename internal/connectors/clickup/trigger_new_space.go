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

package clickup

import (
	"errors"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type TriggerNewSpace struct {
	options *sdk.TriggerInfo
}

func NewTriggerNewSpace() *TriggerNewSpace {
	return &TriggerNewSpace{
		options: &sdk.TriggerInfo{
			Name:        "New Space",
			Description: "triggers workflow when a new space is created",
			RequireAuth: true,
			Auth:        sharedAuth,
			Type:        sdkcore.TriggerTypeCron,
			Input:       map[string]*sdkcore.AutoFormSchema{},
			Settings:    &sdkcore.TriggerSettings{},
			ErrorSettings: &sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (t *TriggerNewSpace) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing clickup auth token")
	}

	return nil, nil
}

func (t *TriggerNewSpace) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return t.Run(ctx)
}

func (t *TriggerNewSpace) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewSpace) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewSpace) GetInfo() *sdk.TriggerInfo {
	return t.options
}
