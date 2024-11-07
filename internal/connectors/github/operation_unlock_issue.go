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

package github

import (
	"errors"
	"fmt"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type unlockIssueOperationProps struct {
	Repository  string `json:"repository"`
	IssueNumber string `json:"issue_number"`
}

type UnlockIssueOperation struct {
	options *sdk.OperationInfo
}

func NewUnlockIssueOperation() *UnlockIssueOperation {
	return &UnlockIssueOperation{
		options: &sdk.OperationInfo{
			Name:        "Unlock Issue",
			Description: "Unlocks an issue",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"repository":   getRepositoryInput(),
				"issue_number": getIssuesInput(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *UnlockIssueOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input, err := sdk.InputToTypeSafely[unlockIssueOperationProps](ctx)
	if err != nil {
		return nil, err
	}

	mutation := fmt.Sprintf(`
		mutation {
		unlockLockable(input: { lockableId: "%s" }) {
			    unlockedRecord{
    				locked
  				}
	    }
	}`, input.IssueNumber)

	response, err := githubGQL(ctx.Auth.AccessToken, mutation)
	if err != nil {
		return nil, errors.New("error making graphQL request")
	}

	issue, ok := response["data"].(map[string]interface{})["unlockLockable"].(map[string]interface{})["unlockedRecord"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to extract issue from response")
	}

	return issue, nil
}

func (c *UnlockIssueOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *UnlockIssueOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
