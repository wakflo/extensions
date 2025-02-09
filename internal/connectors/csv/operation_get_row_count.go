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
	"github.com/gocarina/gocsv"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type getRowCountOperationProps struct {
	Script string `json:"script"`
}

type GetRowCountOperation struct {
	options *sdk.OperationInfo
}

func NewGetRowCountOperation() *GetRowCountOperation {
	return &GetRowCountOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Row Count",
			Description: "retrieves the row count of a csv file",
			Input: map[string]*sdkcore.AutoFormSchema{
				"script": autoform.NewCodeEditorField(sdkcore.CodeLanguageJavascript).
					SetDisplayName("CodeLanguageJavascript Code").
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

func (c *GetRowCountOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	_ = sdk.InputToType[getRowCountOperationProps](ctx)

	res, err := gocsv.MarshalString("gello")
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *GetRowCountOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetRowCountOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
