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

package googledrive

import (
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type CreateNewFileOperation struct {
	options *sdk.OperationInfo
}

func NewCreateNewFileTrigger() *CreateNewFileOperation {
	return &CreateNewFileOperation{
		options: &sdk.OperationInfo{
			Name:        "Create New File",
			Description: "operation creates new google drive file",
			RequireAuth: true,
			Auth:        sharedAuth,
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *CreateNewFileOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	// fake implementation to prove test
	return nil, nil
}

func (c *CreateNewFileOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateNewFileOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
