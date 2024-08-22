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

package jsonconverter

import (
	"fmt"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type convertToJSONOperationProps struct {
	Text string `json:"text"`
}

type ConvertToJSONOperation struct {
	options *sdk.OperationInfo
}

func NewConvertToJSONOperation() *ConvertToJSONOperation {
	return &ConvertToJSONOperation{
		options: &sdk.OperationInfo{
			Name:        "Convert to Json",
			Description: "Returns the text in JSON",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"text": autoform.NewLongTextField().
					SetDisplayName("Text").
					SetDescription("Text to convert").
					SetRequired(true).Build(),
			},
			SampleOutput: map[string]interface{}{
				"text": "Sample text",
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
			RequireAuth: false,
		},
	}
}

func (c *ConvertToJSONOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[convertToJSONOperationProps](ctx)

	result, err := convertTextToJSON(input.Text)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", result)

	return result, nil
}

func (c *ConvertToJSONOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *ConvertToJSONOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
