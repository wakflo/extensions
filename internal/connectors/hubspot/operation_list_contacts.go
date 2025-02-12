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

package hubspot

import (
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
	sdk2 "github.com/wakflo/go-sdk/sdk"
)

type listContactsProps struct {
	Limit string `json:"limit"`
}

type ListContactsOperation struct {
	options *sdk.OperationInfo
}

func NewListContactsOperation() *ListContactsOperation {
	return &ListContactsOperation{
		options: &sdk.OperationInfo{
			Name:        "List Contacts",
			Description: "list contacts",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"limit": autoform.NewShortTextField().
					SetDisplayName("limit").
					SetDescription("").
					SetRequired(false).Build(),
			},
			SampleOutput: map[string]interface{}{},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
			RequireAuth: true,
		},
	}
}

func (c *ListContactsOperation) Run(ctx *sdk.RunContext) (sdk2.JSON, error) {
	_ = sdk.InputToType[listContactsProps](ctx)

	reqURL := "/crm/v3/objects/contacts"

	resp, err := hubspotClient(reqURL, ctx.Auth.AccessToken, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ListContactsOperation) Test(ctx *sdk.RunContext) (sdk2.JSON, error) {
	return c.Run(ctx)
}

func (c *ListContactsOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
