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

package linear

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type TriggerIssueUpdated struct {
	options *sdk.TriggerInfo
}

func NewTriggerIssueUpdated() *TriggerIssueUpdated {
	return &TriggerIssueUpdated{
		options: &sdk.TriggerInfo{
			Name:        "Issue Updated",
			Description: "Triggers workflow when an issue is updated",
			RequireAuth: true,
			Auth:        sharedAuth,
			Type:        sdkcore.TriggerTypeCron,
			Input:       map[string]*sdkcore.AutoFormSchema{},
			Settings:    &sdkcore.TriggerSettings{},
			ErrorSettings: &sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (t *TriggerIssueUpdated) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing linear api key")
	}
	apiKEY := ctx.Auth.Extra["api-key"]

	if !strings.HasPrefix(apiKEY, "lin_api_") {
		return nil, errors.New("invalid Linear API key: must start with 'lin_api_'")
	}

	lastRunTime := ctx.Metadata.LastRun

	var query string
	if lastRunTime == nil {
		query = `{
				issues {
					nodes {
						id
						title
						updatedAt
					}
				}
			}`
	} else {
		query = fmt.Sprintf(`{
				issues(filter: {updatedAt: {gt: "%s"}}) {
					nodes {
						id
						title
						updatedAt
					}
				}
			}`, lastRunTime.UTC().Format(time.RFC3339))
	}

	response, err := MakeGraphQLRequest(apiKEY, query)
	if err != nil {
		log.Fatalf("Error making GraphQL request: %v", err)
	}

	return map[string]interface{}{
		"Result": response,
	}, nil
}

func (t *TriggerIssueUpdated) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return t.Run(ctx)
}

func (t *TriggerIssueUpdated) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerIssueUpdated) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerIssueUpdated) GetInfo() *sdk.TriggerInfo {
	return t.options
}
