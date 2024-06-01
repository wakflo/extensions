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

package starlark

import (
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type RunStarlarkOperation struct {
	options *sdk.OperationInfo
}

func NewRunStarlarkOperation() *RunStarlarkOperation {
	return &RunStarlarkOperation{
		options: &sdk.OperationInfo{
			Name:        "Run Starlark",
			Description: "runs starlark codes",
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *RunStarlarkOperation) Run(ctx *sdk.RunContext) (sdk.Json, error) {
	// todo: rex implement connector
	// fake implementation to prove test
	return map[string]interface{}{
		"data": "some data",
	}, nil
}

func (c *RunStarlarkOperation) Test(ctx *sdk.RunContext) (sdk.Json, error) {
	return c.Run(ctx)
}

func (c *RunStarlarkOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
