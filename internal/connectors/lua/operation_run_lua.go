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

package lua

import (
	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type RunJSOperation struct {
	options *sdk.OperationInfo
}

func NewRunLuaScriptOperation() *RunJSOperation {
	return &RunJSOperation{
		options: &sdk.OperationInfo{
			Name:        "Run Lua Script",
			Description: "runs lua codes",
			Input: map[string]*sdkcore.AutoFormSchema{
				"script": autoform.NewCodeEditorField(sdkcore.Javascript).
					SetDisplayName("Lua Code").
					SetDescription("Enter your lua code").
					SetDefaultValue("function execute(ctx){\n    return {}\n}").
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

func (c *RunJSOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	// todo: rex implement connector
	// fake implementation to prove test
	return map[string]interface{}{
		"data": "some data",
	}, nil
}

func (c *RunJSOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *RunJSOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
