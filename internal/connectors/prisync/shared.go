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

package prisync

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	sdkcore "github.com/wakflo/go-sdk/core"

	"github.com/wakflo/go-sdk/autoform"
)

var sharedAuth = autoform.NewCustomAuthField().
	SetFields(map[string]*sdkcore.AutoFormSchema{
		"api-key": autoform.NewShortTextField().SetDisplayName("Api Key (Required*)").
			SetDescription("The api key used to authenticate prisync.").
			SetRequired(false).
			Build(),
		"api-token": autoform.NewShortTextField().SetDisplayName("Api Token (Required*)").
			SetDescription("The api token used to authenticate prisync.").
			SetRequired(true).
			Build(),
	}).
	Build()

const baseAPI = "https://prisync.com"

func prisyncRequest(apiKey, apiToken, reqURL, method string, formData map[string]string) (interface{}, error) {
	fullRequest := baseAPI + reqURL
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	for key, val := range formData {
		if err := writer.WriteField(key, val); err != nil {
			return nil, err
		}
	}

	writer.Close()

	req, err := http.NewRequest(method, fullRequest, &body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Apikey", apiKey)
	req.Header.Add("Apitoken", apiToken)

	client := &http.Client{}
	res, errs := client.Do(req)
	if errs != nil {
		return nil, errs
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response interface{}
	if newErrs := json.Unmarshal(respBody, &response); newErrs != nil {
		return nil, errors.New("error parsing response")
	}

	return response, nil
}

func prisyncBatchRequest(apiKey, apiToken, reqURL string, products []map[string]string, cancelOnPackageLimitExceeding bool) (interface{}, error) {
	fullRequest := baseAPI + reqURL

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	for i, product := range products {
		for key, value := range product {
			fieldName := fmt.Sprintf("product%d[%s]", i, key)
			if err := writer.WriteField(fieldName, value); err != nil {
				return nil, err
			}
		}
	}

	err := writer.WriteField("cancelOnPackageLimitExceeding", "false")
	if err != nil {
		return nil, err
	}

	writer.Close()

	req, err := http.NewRequest(http.MethodPost, fullRequest, &body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Apikey", apiKey)
	req.Header.Add("Apitoken", apiToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response interface{}
	if errs := json.Unmarshal(respBody, &response); errs != nil {
		return nil, errors.New("error parsing response")
	}

	return response, nil
}
