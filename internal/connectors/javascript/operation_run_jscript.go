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
	"errors"
	"time"

	"github.com/dop251/goja"
	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type runJSOperationProps struct {
	Script string `json:"script"`
}

type RunJSOperation struct {
	options *sdk.OperationInfo
}

func NewRunJSOperation() *RunJSOperation {
	return &RunJSOperation{
		options: &sdk.OperationInfo{
			Name:        "Run JavaScript",
			Description: "runs javascript codes",
			Input: map[string]*sdkcore.AutoFormSchema{
				"script": autoform.NewCodeEditorField(sdkcore.Javascript).
					SetDisplayName("Javascript Code").
					SetDescription("Enter your javascript code").
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

func (c *RunJSOperation) Run(ctx *sdk.RunContext) (sdk.Json, error) {
	input := sdk.InputToType[runJSOperationProps](ctx)

	vm := goja.New()
	// map field names to JS
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	// we want to stop processing if it takes more that 500ms
	time.AfterFunc(200*time.Millisecond, func() {
		vm.Interrupt("ran out of time")
	})

	// Run Script
	_, err := vm.RunScript(ctx.Step.Name, input.Script)
	if err != nil {
		return nil, err
	}

	// assert if function exists
	execute, ok := goja.AssertFunction(vm.Get("execute"))
	if !ok {
		return nil, errors.New("missing execute function")
	}

	arg := map[string]interface{}{
		"step":  ctx.Step,
		"steps": ctx.State.Steps,
		"auth":  ctx.Auth,
		"input": input,
	}
	res, err := execute(goja.Undefined(), vm.ToValue(arg))
	if err != nil {
		return nil, err
	}

	return res.Export(), nil
}

func (c *RunJSOperation) Test(ctx *sdk.RunContext) (sdk.Json, error) {
	return c.Run(ctx)
}

func (c *RunJSOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
