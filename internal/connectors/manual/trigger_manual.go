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

package javascript

import (
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type ManualTrigger struct {
	options *sdk.TriggerInfo
}

func NewManualTrigger() *ManualTrigger {
	return &ManualTrigger{
		options: &sdk.TriggerInfo{
			Name:         "Trigger",
			Description:  "manual trigger",
			Input:        map[string]*sdkcore.AutoFormSchema{},
			SampleOutput: map[string]interface{}{},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
			Settings: &sdkcore.TriggerSettings{
				Type: sdkcore.Manual,
			},
			RequireAuth: false,
		},
	}
}

func (c *ManualTrigger) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	return map[string]interface{}{}, nil
}

func (c *ManualTrigger) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *ManualTrigger) GetInfo() *sdk.TriggerInfo {
	return c.options
}

func (c *ManualTrigger) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (c *ManualTrigger) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}
