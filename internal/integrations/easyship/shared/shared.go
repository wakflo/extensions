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
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	fastshot "github.com/opus-domini/fast-shot"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

const baseURL = "https://api.easyship.com"

func NewEasyShipAPIClient(baseURL, apiKey string) *http.Client {
	return &http.Client{}
}

func PostRequest(endpoint, apiKey string, labelData map[string]interface{}) (map[string]interface{}, error) {
	client := NewEasyShipAPIClient(baseURL, apiKey)

	jsonData, err := json.Marshal(labelData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var result map[string]interface{}
	if newErr := json.Unmarshal(body, &result); newErr != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", newErr)
	}

	return result, nil
}

func RegisterCourierProps(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getCouriers := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		// Get API key from auth context
		apiKey, ok := authCtx.Extra["api-key"]
		if !ok {
			return nil, errors.New("API key not found in auth context")
		}

		// Create request to get couriers
		qu := fastshot.NewClient(baseURL).
			Auth().BearerToken(apiKey).
			Header().
			AddAccept("application/json").
			Build().GET("/2024-09/couriers")

		rsp, err := qu.Send()
		if err != nil {
			return nil, err
		}

		if rsp.Status().IsError() {
			return nil, errors.New(rsp.Status().Text())
		}
		defer rsp.Raw().Body.Close()

		bytes, err := io.ReadAll(rsp.Raw().Body)
		if err != nil {
			return nil, err
		}

		// Define response structure based on Easyship API
		var response struct {
			Couriers []struct {
				ID           string `json:"id"`
				UmbrellaName string `json:"umbrella_name"`
			} `json:"couriers"`
		}

		err = json.Unmarshal(bytes, &response)
		if err != nil {
			return nil, err
		}

		var options []map[string]interface{}
		for _, courier := range response.Couriers {
			options = append(options, map[string]interface{}{
				"id":   courier.ID,
				"name": courier.UmbrellaName,
			})
		}

		return ctx.Respond(options, len(options))
	}

	return form.SelectField("courier-id", "Courier ID").
		Placeholder("Select a courier").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getCouriers)).
				WithSearchSupport().
				End().
				GetDynamicSource(),
		).
		HelpText("Select a courier service for shipping")
}
