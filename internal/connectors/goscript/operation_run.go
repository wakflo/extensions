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

package goscript

import (
	"fmt"

	"github.com/open2b/scriggo"
	"github.com/open2b/scriggo/native"
	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type runGoScriptOperationProps struct {
	Script string `json:"script"`
}

var packages native.Packages

type RunGoScriptOperation struct {
	options *sdk.OperationInfo
	vmopts  *scriggo.BuildOptions
}

func NewRunGoScriptOperation() *RunGoScriptOperation {
	return &RunGoScriptOperation{
		options: &sdk.OperationInfo{
			Name:        "Run Go",
			Description: "runs golang codes",
			Input: map[string]*sdkcore.AutoFormSchema{
				"script": autoform.NewCodeEditorField(sdkcore.GoLang).
					SetDisplayName("GoLang Code").
					SetDescription("Enter your go code").
					SetDefaultValue("func execute(ctx any) map[string]interface{} {\n\treturn map[string]interface{}{}\n}").
					SetRequired(true).Build(),
			},
			SampleOutput: map[string]interface{}{},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
			RequireAuth: false,
		},
		vmopts: &scriggo.BuildOptions{Packages: packages},
	}
}

func (c *RunGoScriptOperation) Run(ctx *sdk.RunContext) (sdk.Json, error) {
	input := sdk.InputToType[runGoScriptOperationProps](ctx)

	src := fmt.Sprintf(`package main

func main() {
	execute(nil)
}
%v
`, input.Script)
	// Create a file system with the file of the program to run.
	fsys := scriggo.Files{"main.go": []byte(src)}

	//arg := map[string]interface{}{
	//	"step":     ctx.Step,
	//	"auth":     ctx.Auth,
	//	"workflow": ctx.Workflow,
	//}

	// Build the program.
	program, err := scriggo.Build(fsys, c.vmopts)
	if err != nil {
		return nil, err
	}

	// Run the program.
	err = program.Run(nil)
	if err != nil {
		return nil, err
	}

	return nil, err
}

func (c *RunGoScriptOperation) Test(ctx *sdk.RunContext) (sdk.Json, error) {
	return c.Run(ctx)
}

func (c *RunGoScriptOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
