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

type deleteCardOperationProps struct {
	CardID string `json:"cardId"`
}

type DeleteCardOperation struct {
	options *sdk.OperationInfo
}

func NewDeleteCardOperation() *DeleteCardOperation {
	return &DeleteCardOperation{
		options: &sdk.OperationInfo{
			Name:        "Delete a Card",
			Description: "delete a card",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"cardId": autoform.NewShortTextField().
					SetDisplayName("Card ID").
					SetDescription("The id of the card to delete").
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

func (c *DeleteCardOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing trello api credentials")
	}
	apiKey := ctx.Auth.Extra["api-key"]
	apiToken := ctx.Auth.Extra["api-token"]
	input := sdk.InputToType[deleteCardOperationProps](ctx)
	fullURL := fmt.Sprintf("%s/cards/%s?key=%s&token=%s", baseURL, input.CardID, apiKey, apiToken)

	_, err := trelloRequest(http.MethodDelete, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	return map[string]interface{}{
		"Result": "Card deleted Successfully",
	}, nil
}

func (c *DeleteCardOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *DeleteCardOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
