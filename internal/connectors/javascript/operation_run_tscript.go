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
	"strings"

	"github.com/clarkmcc/go-typescript"
	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type runTSOperationProps struct {
	Script string `json:"script"`
}

type RunTSOperation struct {
	options *sdk.OperationInfo
}

func NewRunTSOperation() sdk.IOperation {
	return &RunTSOperation{
		options: &sdk.OperationInfo{
			Name:        "Run TypeScript",
			Description: "runs typescript codes",
			Input: map[string]*sdkcore.AutoFormSchema{
				"script": autoform.NewCodeEditorField(sdkcore.Typescript).
					SetDisplayName("TypeScript Code").
					SetDescription("Enter your typescript code").
					SetDefaultValue("function execute(ctx: any){\n    return {}\n}").
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

func (c *RunTSOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[runTSOperationProps](ctx)

	res, err := typescript.Evaluate(
		strings.NewReader(input.Script),
		typescript.WithTranspile(),
	)
	if err != nil {
		return nil, err
	}

	return res.Export(), nil
}

func (c *RunTSOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *RunTSOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
