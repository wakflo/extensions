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

package delay

import (
	"time"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
	sdk2 "github.com/wakflo/go-sdk/sdk"
)

type delayUntilOperationOutput struct {
	DelayUntil time.Time `json:"delayInMs"`
	Success    bool      `json:"success"`
}

type delayUntilOperationProps struct {
	DelayUntil string `json:"delayUntil"`
}

type DelayUntilOperation struct {
	options *sdk.OperationInfo
}

func NewDelayUntilOperation() *DelayUntilOperation {
	return &DelayUntilOperation{
		options: &sdk.OperationInfo{
			Name:        "Delay Until",
			Description: "Delays the execution of the next action until a given timestamp",
			Input: map[string]*sdkcore.AutoFormSchema{
				"delayUntil": autoform.NewDateTimeField().
					SetDisplayName("Date and Time").
					SetDescription("Specifies the date and time until which the execution of the next action should be delayed. It supports multiple formats, including ISO format.").
					SetRequired(true).Build(),
			},
			SampleOutput: map[string]interface{}{},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
			RequireAuth: false,
		},
	}
}

func (c *DelayUntilOperation) Run(ctx *sdk.RunContext) (sdk2.JSON, error) {
	input := sdk.InputToType[delayUntilOperationProps](ctx)

	delayTill, err := time.Parse(time.RFC3339, input.DelayUntil)
	if err != nil {
		return nil, err
	}

	delayInMs := delayTill.UnixMilli() - time.Now().UnixMilli()
	if ctx.ExecutionType == sdkcore.Resume || delayInMs <= 0 {
		return &delayUntilOperationOutput{
			DelayUntil: delayTill,
			Success:    true,
		}, nil
	}

	if delayInMs > 1*60*1000 {
		futureTime := time.Now().Add(time.Millisecond * time.Duration(delayInMs))
		return ctx.PauseExecution(sdk.PauseMetadata{
			Type:     sdkcore.DelayPause,
			ResumeAt: &futureTime,
		})
	}

	time.Sleep(time.Millisecond * time.Duration(delayInMs))
	return &delayUntilOperationOutput{
		DelayUntil: delayTill,
		Success:    true,
	}, nil
}

func (c *DelayUntilOperation) Test(ctx *sdk.RunContext) (sdk2.JSON, error) {
	return c.Run(ctx)
}

func (c *DelayUntilOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
