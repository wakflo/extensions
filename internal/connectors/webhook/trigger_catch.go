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

package webhook

import (
	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type runJSOperationProps struct {
	Script string `json:"script"`
}

type CatchTrigger struct {
	options *sdk.TriggerInfo
}

func NewCatchTrigger() sdk.ITrigger {
	return &CatchTrigger{
		options: &sdk.TriggerInfo{
			Name:        "Catch",
			Description: "Receive incoming HTTP/webhooks using any HTTP method such as GET, POST, PUT, DELETE, etc. ",
			Input: map[string]*sdkcore.AutoFormSchema{
				"Public": autoform.NewBooleanField().
					SetDisplayName("Public").
					SetDescription("Whether the webhook endpoint does not require authentication and is public (a true or false value).").
					SetRequired(false).Build(),
			},
			Type:         sdkcore.TriggerTypeWebhook,
			SampleOutput: map[string]interface{}{},
			Settings:     &sdkcore.TriggerSettings{},
			ErrorSettings: &sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
			RequireAuth:   false,
			Documentation: &triggerCatchDocs,
			HelpText:      &helpText,
		},
	}
}

func (c *CatchTrigger) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input, err := sdk.InputToTypeSafely[runJSOperationProps](ctx)
	return input, err
}

func (c *CatchTrigger) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CatchTrigger) GetInfo() *sdk.TriggerInfo {
	return c.options
}

func (c *CatchTrigger) OnEnabled(ctx *sdk.RunContext) error {
	// TODO implement me
	panic("implement me")
}

func (c *CatchTrigger) OnDisabled(ctx *sdk.RunContext) error {
	// TODO implement me
	panic("implement me")
}
