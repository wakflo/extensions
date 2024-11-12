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

package zendesk

import (
	"errors"
	"fmt"
	"net/http"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type GetGroupsOperation struct {
	options *sdk.OperationInfo
}

func NewGetGroupsOperation() *GetGroupsOperation {
	return &GetGroupsOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Groups",
			Description: "get groups",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input:       map[string]*sdkcore.AutoFormSchema{},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *GetGroupsOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-token"] == "" || ctx.Auth.Extra["email"] == "" || ctx.Auth.Extra["subdomain"] == "" {
		return nil, errors.New("missing zendesk api credentials")
	}

	email := ctx.Auth.Extra["email"]
	apiToken := ctx.Auth.Extra["api-token"]
	subdomain := "https://" + ctx.Auth.Extra["subdomain"] + ".zendesk.com/api/v2"

	fullURL := subdomain + "/groups.json"

	response, err := zendeskRequest(http.MethodGet, fullURL, email, apiToken, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	return response, nil
}

func (c *GetGroupsOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetGroupsOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
