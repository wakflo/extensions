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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	sdk "github.com/wakflo/go-sdk/connector"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var sharedAuth = autoform.NewCustomAuthField().
	SetFields(map[string]*sdkcore.AutoFormSchema{
		"api-key": autoform.NewShortTextField().SetDisplayName("Api Key (Required)").
			SetDescription("Your trello api key").
			SetRequired(true).
			Build(),
		"api-token": autoform.NewLongTextField().SetDisplayName("Api Token (Required)").
			SetDescription("Your Trello API Token. Click **manually generate a Token** next to the API key field").
			SetRequired(true).
			Build(),
	}).
	Build()

const baseURL = "https://api.trello.com/1"

func trelloRequest(method, fullURL string, request []byte) (interface{}, error) {
	req, err := http.NewRequest(method, fullURL, bytes.NewBuffer(request))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	var result interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return result, nil
}

func getBoardsInput() *sdkcore.AutoFormSchema {
	getBoards := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
		endpoint := "/members/me/boards"
		fullURL := fmt.Sprintf("%s%s?key=%s&token=%s", baseURL, endpoint, ctx.Auth.Extra["api-key"], ctx.Auth.Extra["api-token"])

		req, err := http.NewRequest(http.MethodGet, fullURL, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/json")

		client := &http.Client{}
		rsp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer rsp.Body.Close()

		newBytes, err := io.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}

		var boards []Board
		err = json.Unmarshal(newBytes, &boards)
		if err != nil {
			return nil, err
		}

		return boards, nil
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Boards").
		SetDescription("Select a board").
		SetDynamicOptions(&getBoards).
		SetRequired(true).Build()
}

func getBoardListsInput() *sdkcore.AutoFormSchema {
	getBoardLists := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
		input := sdk.DynamicInputToType[struct {
			BoardID string `json:"board_id"`
			ListID  string `json:"idList"`
		}](ctx)

		endpoint := fmt.Sprintf("/boards/%s/lists", input.BoardID)
		fullURL := fmt.Sprintf("%s%s?key=%s&token=%s", baseURL, endpoint, ctx.Auth.Extra["api-key"], ctx.Auth.Extra["api-token"])

		req, err := http.NewRequest(http.MethodGet, fullURL, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/json")

		client := &http.Client{}
		rsp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer rsp.Body.Close()

		newBytes, err := io.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}

		var lists []BoardList
		err = json.Unmarshal(newBytes, &lists)
		if err != nil {
			return nil, err
		}

		return lists, nil
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Lists").
		SetDescription("Select a list").
		SetDynamicOptions(&getBoardLists).
		SetRequired(true).Build()
}

// func getCardsInput() *sdkcore.AutoFormSchema {
//	getCards := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
//		input := sdk.DynamicInputToType[struct {
//			BoardID string `json:"board_id"`
//			ListID  string `json:"idList"`
//		}](ctx)
//
//		fullURL := fmt.Sprintf("%s/lists/%s/cards?key=%s&token=%s", baseURL, input.ListID, ctx.Auth.Extra["api-key"], ctx.Auth.Extra["api-token"])
//
//		req, err := http.NewRequest(http.MethodGet, fullURL, nil)
//		if err != nil {
//			return nil, err
//		}
//
//		req.Header.Set("Accept", "application/json")
//
//		client := &http.Client{}
//		rsp, err := client.Do(req)
//		if err != nil {
//			return nil, err
//		}
//		defer rsp.Body.Close()
//
//		newBytes, err := io.ReadAll(rsp.Body)
//		if err != nil {
//			return nil, err
//		}
//
//		var lists []Cards
//		err = json.Unmarshal(newBytes, &lists)
//		if err != nil {
//			return nil, err
//		}
//
//		return lists, nil
//	}
//
//	return autoform.NewDynamicField(sdkcore.String).
//		SetDisplayName("Cards").
//		SetDescription("Select a list").
//		SetDynamicOptions(&getCards).
//		SetRequired(true).Build()
// }

var trelloCardPosition = []*sdkcore.AutoFormSchema{
	{Const: "top", Title: "Top"},
	{Const: "bottom", Title: "Bottom"},
}
