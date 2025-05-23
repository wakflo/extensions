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
	"io"
	"net/http"
	"net/url"

	"github.com/juicycleff/smartform/v1"
)

var (
	form = smartform.NewAuthForm("stripe-auth", "Stripe API Authentication", smartform.AuthStrategyCustom)

	_ = form.TextField("api-key", "Secret API Key (Required)").
		Required(true).
		HelpText("Secret key acquired from your Stripe dashboard")

	StripeSharedAuth = form.Build()
)

const baseURL = "https://api.stripe.com"

func StripeClient(apiKey, url, httpType string, payload []byte, params url.Values) (map[string]interface{}, error) {
	fullURL := baseURL + url

	req, err := http.NewRequest(httpType, fullURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	if len(params) > 0 {
		req.URL.RawQuery = params.Encode()
	}

	req.SetBasicAuth(apiKey, "")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, errors.New("error unmarshalling response")
	}

	return result, nil
}
