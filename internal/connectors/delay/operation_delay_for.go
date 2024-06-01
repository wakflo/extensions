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
)

type delayForOperationProps struct {
	Unit     TimeUnit `json:"unit"`
	DelayFor int      `json:"delayFor"`
}

type delayForOperationOutput struct {
	DelayInMs int  `json:"delayInMs"`
	Success   bool `json:"success"`
}

type DelayForOperation struct {
	options *sdk.OperationInfo
}

func NewDelayForOperation() *DelayForOperation {
	return &DelayForOperation{
		options: &sdk.OperationInfo{
			Name:        "Delay For",
			Description: "Delays the execution of the next action for a given duration",
			Input: map[string]*sdkcore.AutoFormSchema{
				"unit": autoform.NewSelectField().
					SetDisplayName("Unit").
					SetDescription("The unit of time to delay the execution of the next operation").
					SetOptions(timeUnitWithLabels).
					SetDefaultValue(Seconds).
					SetRequired(true).
					Build(),
				"delayFor": autoform.NewNumberField().
					SetDisplayName("Delay Amount").
					SetDescription("The number of units to delay the execution of the next operation").
					SetRequired(true).
					Build(),
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

func (c *DelayForOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[delayForOperationProps](ctx)

	unit := Seconds
	if input.Unit != "" {
		unit = input.Unit
	}

	var delayInMs int
	switch unit {
	case Seconds:
		delayInMs = input.DelayFor * 1000
	case Minutes:
		delayInMs = input.DelayFor * 60 * 1000
	case Hours:
		delayInMs = input.DelayFor * 60 * 60 * 1000
	case Days:
		delayInMs = input.DelayFor * 24 * 60 * 60 * 1000
	}

	if ctx.ExecutionType == sdkcore.Resume {
		return &delayForOperationOutput{
			DelayInMs: delayInMs,
			Success:   true,
		}, nil
	}

	if delayInMs > 1*60*1000 {
		futureTime := time.Now().Add(time.Millisecond * time.Duration(delayInMs))
		return ctx.PauseExecution(sdk.PauseMetadata{
			Type:     sdkcore.DelayPause,
			ResumeAt: &futureTime,
		})
	}

	// sleep for a bit
	time.Sleep(time.Millisecond * time.Duration(delayInMs))

	return &delayForOperationOutput{
		DelayInMs: delayInMs,
		Success:   true,
	}, nil
}

func (c *DelayForOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *DelayForOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
