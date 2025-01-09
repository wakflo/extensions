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

package monday

import (
	"errors"
	"fmt"
	"strings"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createUpdateOperationProps struct {
	ItemID string `json:"item_id,omitempty"`
	Body   string `json:"body,omitempty"`
}

type CreateUpdateOperation struct {
	options *sdk.OperationInfo
}

func NewCreateUpdateOperation() *CreateUpdateOperation {
	return &CreateUpdateOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Update",
			Description: "Adds an update to an item",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"item_id": autoform.NewShortTextField().
					SetDisplayName("Item ID").
					SetDescription("Item ID").
					SetRequired(true).
					Build(),
				"body": autoform.NewLongTextField().
					SetDisplayName("Body").
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

func (c *CreateUpdateOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing monday.com auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input, err := sdk.InputToTypeSafely[createUpdateOperationProps](ctx)
	if err != nil {
		return nil, err
	}

	fields := make(map[string]string)
	fields["item_id"] = fmt.Sprintf(`"%s"`, input.ItemID)
	fields["body"] = fmt.Sprintf(`"%s"`, input.Body)

	fieldStrings := make([]string, 0, len(fields))
	for key, value := range fields {
		fieldStrings = append(fieldStrings, fmt.Sprintf("%s: %s", key, value))
	}

	query := fmt.Sprintf(`
		mutation {
  create_update (%s) {
    	id
		body
  }
}`, strings.Join(fieldStrings, "\n"))

	response, err := mondayClient(accessToken, query)
	if err != nil {
		return nil, err
		// return nil, errors.New("request not successful")
	}

	update, ok := response["data"].(map[string]interface{})["create_update"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to extract item from response")
	}

	return update, nil
}

func (c *CreateUpdateOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateUpdateOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
