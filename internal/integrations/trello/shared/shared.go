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

package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"

	"github.com/wakflo/go-sdk/v2"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

var (
	form = smartform.NewAuthForm("trello-auth", "Trello API Authentication", smartform.AuthStrategyCustom)

	_ = form.TextField("api-key", "Api Key(Required)").
		Required(true).
		HelpText("Your trello api key ")

	_ = form.TextField("api-token", "Api Token (Required)").
		Required(true).
		HelpText("Your Trello API Token. Click **manually generate a Token** next to the API key field.")

	TrelloSharedAuth = form.Build()
)

const BaseURL = "https://api.trello.com/1"

func TrelloRequest(method, fullURL string, request []byte) (interface{}, error) {
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

func RegisterBoardsProp(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getBoards := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}
		endpoint := "/members/me/boards"
		fullURL := fmt.Sprintf("%s%s?key=%s&token=%s", BaseURL, endpoint, authCtx.Extra["api-key"], authCtx.Extra["api-token"])

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

		return ctx.Respond(boards, len(boards))
	}

	return form.SelectField("boards", "Boards").
		Placeholder("Select a board.").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getBoards)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText("Select a board")
}

func RegisterBoardListsProp(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getBoardLists := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		input := sdk.DynamicInputToType[struct {
			BoardID string `json:"boards"`
			ListID  string `json:"idList"`
		}](ctx)

		endpoint := fmt.Sprintf("/boards/%s/lists", input.BoardID)
		fullURL := fmt.Sprintf("%s%s?key=%s&token=%s", BaseURL, endpoint, authCtx.Extra["api-key"], authCtx.Extra["api-token"])

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

		return ctx.Respond(lists, len(lists))
	}

	return form.SelectField("list", "Lists").
		Placeholder("Select a list.").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getBoardLists)).
				WithSearchSupport().
				WithPagination(10).
				End().
				RefreshOn("boards").
				GetDynamicSource(),
		).
		HelpText("Select a list")
}

func RegisterCardsProp(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getCards := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		input := sdk.DynamicInputToType[struct {
			BoardID string `json:"boards"`
			ListID  string `json:"list"`
		}](ctx)

		endpoint := fmt.Sprintf("/lists/%s/cards", input.ListID)
		fullURL := fmt.Sprintf("%s%s?key=%s&token=%s", BaseURL, endpoint, authCtx.Extra["api-key"], authCtx.Extra["api-token"])

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

		var cards []Card

		err = json.Unmarshal(newBytes, &cards)
		if err != nil {
			return nil, err
		}

		return ctx.Respond(cards, len(cards))
	}

	return form.SelectField("cardId", "Cards").
		Placeholder("Select a card.").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getCards)).
				WithSearchSupport().
				WithPagination(10).
				End().
				RefreshOn("list").
				GetDynamicSource(),
		).
		HelpText("Select a card")
}
