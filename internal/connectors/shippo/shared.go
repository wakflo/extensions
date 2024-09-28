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

package shippo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gookit/goutil/arrutil"
	fastshot "github.com/opus-domini/fast-shot"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var sharedAuth = autoform.NewCustomAuthField().
	SetFields(map[string]*sdkcore.AutoFormSchema{
		"api-key": autoform.NewShortTextField().SetDisplayName("API Key").
			SetDescription("API Application Key").
			SetRequired(true).
			Build(),
	}).
	Build()

const baseURL = "https://api.goshippo.com"

func NewShippoAPIClient(baseURL, apiKey string) *http.Client {
	return &http.Client{}
}

func CreateAShipment(endpoint, apiKey string, shipmentData map[string]interface{}) (interface{}, error) {
	client := NewShippoAPIClient(baseURL, apiKey)

	jsonData, err := json.Marshal(shipmentData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "ShippoToken "+apiKey)
	req.Header.Set("Shippo-Api-Version", "2018-02-08")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var result interface{}
	if newErr := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", newErr)
	}

	return result, nil
}

type Country struct {
	Code string `json:"alpha2Code"`
	Name string `json:"name"`
}

func getCountriesInput() *sdkcore.AutoFormSchema {
	getCountries := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
		qu := fastshot.NewClient("https://restcountries.com").
			Auth().BearerToken(ctx.Auth.AccessToken).
			Header().
			AddAccept("application/json").
			Build().GET("/v2/all")

		rsp, err := qu.Send()
		if err != nil {
			return nil, err
		}

		if rsp.Status().IsError() {
			return nil, errors.New(rsp.Status().Text())
		}
		defer rsp.Raw().Body.Close()

		bytes, err := io.ReadAll(rsp.Raw().Body) //nolint:bodyclose
		if err != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		var response []Country
		err = json.Unmarshal(bytes, &response)
		if err != nil {
			return nil, err
		}

		country := response

		return arrutil.Map[Country, map[string]any](country, func(input Country) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.Code,
				"name": input.Name,
			}, true
		}), nil
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Country").
		SetDescription("select a country").
		SetDynamicOptions(&getCountries).
		SetRequired(false).Build()
}

var distanceUnit = []*sdkcore.AutoFormSchema{
	{Const: "cm", Title: " Centimeters"},
	{Const: "in", Title: "Inches"},
	{Const: "ft", Title: "Feet"},
	{Const: "mm", Title: "Millimeters"},
	{Const: "m", Title: "Meters"},
	{Const: "yd", Title: "Yards"},
}

var massUnit = []*sdkcore.AutoFormSchema{
	{Const: "g", Title: "Grams"},
	{Const: "kg", Title: "Kilogram"},
	{Const: "lb", Title: "Pounds"},
	{Const: "oz", Title: "Ounces"},
}
