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

package trello

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createListOperationProps struct {
	IDBoard  string `json:"idBoard"`
	Name     string `json:"name"`
	Position string `json:"position"`
}

type CreateListOperation struct {
	options *sdk.OperationInfo
}

func NewCreateListOperation() *CreateListOperation {
	return &CreateListOperation{
		options: &sdk.OperationInfo{
			Name:        "Create List",
			Description: "Create a new card on a specific board",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"idBoard": getBoardsInput(),
				"name": autoform.NewShortTextField().
					SetDisplayName("List Name").
					SetDescription("The name of the list to create").
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

func (c *CreateListOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing trello api credentials")
	}
	apiKey := ctx.Auth.Extra["api-key"]
	apiToken := ctx.Auth.Extra["api-token"]
	input := sdk.InputToType[createListOperationProps](ctx)
	fullURL := fmt.Sprintf("%s/lists?name=%s&idBoard=%s&key=%s&token=%s", baseURL, input.Name, input.IDBoard, apiKey, apiToken)

	response, err := trelloRequest(http.MethodPost, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	return response, nil
}

func (c *CreateListOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateListOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
