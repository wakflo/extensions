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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createCardOperationProps struct {
	BoardID     string `json:"board_id"`
	ListID      string `json:"idList"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Position    string `json:"position"`
}

type CreateCardOperation struct {
	options *sdk.OperationInfo
}

func NewCreateCardOperation() *CreateCardOperation {
	return &CreateCardOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Card",
			Description: "Create a new card on a specific board and list.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"board_id": getBoardsInput(),
				"idList":   getBoardListsInput(),
				"name": autoform.NewShortTextField().
					SetDisplayName("Card Name").
					SetDescription("The name of the card to create").
					SetRequired(true).
					Build(),
				"description": autoform.NewLongTextField().
					SetDisplayName("Description").
					SetDescription("The description of the card to create").
					SetRequired(false).
					Build(),
				"position": autoform.NewSelectField().
					SetDisplayName("Position").
					SetDescription("Place the card on top or bottom of the list").
					SetOptions(trelloCardPosition).
					SetRequired(false).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *CreateCardOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing trello api credentials")
	}
	apiKey := ctx.Auth.Extra["api-key"]
	apiToken := ctx.Auth.Extra["api-token"]
	input := sdk.InputToType[createCardOperationProps](ctx)
	fullURL := fmt.Sprintf("%s/cards?idList=%s&key=%s&token=%s", baseURL, input.ListID, apiKey, apiToken)

	payload := CardRequest{
		Name: input.Name,
		Desc: input.Description,
		Pos:  input.Position,
	}

	payloadBytes, _ := json.Marshal(payload)

	response, err := trelloRequest(http.MethodPost, fullURL, payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	return response, nil
}

func (c *CreateCardOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateCardOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
