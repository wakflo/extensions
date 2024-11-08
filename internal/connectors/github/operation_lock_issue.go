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

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type lockIssueOperationProps struct {
	Repository  string `json:"repository"`
	LockReason  string `json:"lock_reason"`
	IssueNumber string `json:"issue_number"`
}

type LockIssueOperation struct {
	options *sdk.OperationInfo
}

func NewLockIssueOperation() *LockIssueOperation {
	return &LockIssueOperation{
		options: &sdk.OperationInfo{
			Name:        "Lock Issue",
			Description: "Locks the specified issue",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"repository":   getRepositoryInput(),
				"issue_number": getIssuesInput(),
				"lock_reason": autoform.NewSelectField().
					SetDisplayName("Lock Reason").
					SetDescription("The reason for locking the issue").
					SetOptions(lockIssueReason).
					SetRequired(true).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *LockIssueOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input, err := sdk.InputToTypeSafely[lockIssueOperationProps](ctx)
	if err != nil {
		return nil, err
	}

	mutation := fmt.Sprintf(`
		mutation {
			lockLockable(input: { lockableId: "%s", lockReason: %s }) {
			lockedRecord {
				locked
        		activeLockReason
			}
		}
	}`, input.IssueNumber, input.LockReason)

	response, err := githubGQL(ctx.Auth.AccessToken, mutation)
	if err != nil {
		return nil, errors.New("error making graphQL request")
	}

	issue, ok := response["data"].(map[string]interface{})["lockLockable"].(map[string]interface{})["lockedRecord"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to extract issue from response")
	}

	return issue, nil
}

func (c *LockIssueOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *LockIssueOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
