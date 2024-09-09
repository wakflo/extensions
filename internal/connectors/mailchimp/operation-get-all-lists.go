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

package mailchimp

import (
	"errors"
	"log"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type GetAllListOperation struct {
	options *sdk.OperationInfo
}

func NewGetAllListOperation() sdk.IOperation {
	return &GetAllListOperation{
		options: &sdk.OperationInfo{
			Name:        "Get all lists ",
			Description: "Get all available lists",
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

func (c *GetAllListOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing mailchimp auth token")
	}
	accessToken := ctx.Auth.AccessToken

	dc, err := getMailChimpServerPrefix(accessToken)
	if err != nil {
		log.Fatalf("Error getting MailChimp server prefix: %v", err)
	}

	var result interface{}
	result, err = fetchMailchimpLists(accessToken, dc)
	if err != nil {
		return nil, err
	}

	return sdk.JSON(map[string]interface{}{
		"result": result,
	}), err
}

func (c *GetAllListOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetAllListOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
