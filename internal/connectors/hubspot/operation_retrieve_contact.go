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
	"fmt"
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type retrieveContactProps struct {
	Email string `json:"email"`
}

type RetrieveContactOperation struct {
	options *sdk.OperationInfo
}

func NewRetrieveContactOperation() *RetrieveContactOperation {
	return &RetrieveContactOperation{
		options: &sdk.OperationInfo{
			Name:        "Retrieve a Contact by email",
			Description: "retrieve a contact",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"email": autoform.NewShortTextField().
					SetDisplayName("Contact's email").
					SetDescription("contact's email").
					SetRequired(true).Build(),
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

func (c *RetrieveContactOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[retrieveContactProps](ctx)

	reqURL := fmt.Sprintf("https://api.hubapi.com/crm/v3/objects/contacts/%s?idProperty=email", input.Email)

	resp, err := hubspotClient(reqURL, ctx.Auth.AccessToken, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *RetrieveContactOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *RetrieveContactOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
