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
	"mime/multipart"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

var (
	form = smartform.NewAuthForm("prisync-auth", "Prisync API Authentication", smartform.AuthStrategyCustom)

	_ = form.TextField("api-key", "API Key (Required*)").
		Required(true).
		HelpText("The api key used to authenticate prisync.")

	_ = form.TextField("api-token", "API Token (Required*)").
		Required(true).
		HelpText("The api token used to authenticate prisync.")

	PrisyncSharedAuth = form.Build()
)

const baseAPI = "https://prisync.com"

func PrisyncRequest(apiKey, apiToken, reqURL, method string, formData map[string]string) (interface{}, error) {
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

func PrisyncBatchRequest(apiKey, apiToken, reqURL string, products []map[string]string, cancelOnPackageLimitExceeding bool) (interface{}, error) {
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

func GetProductProp(id string, title string, desc string, required bool, form *smartform.FormBuilder) *smartform.FieldBuilder {
	listPrisyncProducts := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		// Get API credentials from context
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		apiKey := authCtx.Extra["api-key"]
		apiToken := authCtx.Extra["api-token"]

		// Prisync API endpoint for fetching products
		endpoint := "/api/v2/list/product/startFrom/0"

		// Use the Prisync client to fetch products
		response, err := PrisyncRequest(apiKey, apiToken, endpoint, http.MethodGet, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch Prisync products: %v", err)
		}

		// Process the response
		var options []map[string]interface{}

		// Parse the response based on Prisync API structure
		if responseMap, ok := response.(map[string]interface{}); ok {
			// Extract the results array
			if results, ok := responseMap["results"].([]interface{}); ok {
				for _, item := range results {
					productMap, ok := item.(map[string]interface{})
					if !ok {
						continue
					}

					// Extract product properties
					id, idOk := productMap["id"]
					name, nameOk := productMap["name"].(string)

					if !idOk || !nameOk {
						continue
					}

					// Convert ID to string based on type
					var idStr string
					switch v := id.(type) {
					case float64:
						idStr = fmt.Sprintf("%.0f", v)
					case string:
						idStr = v
					case json.Number:
						idStr = string(v)
					default:
						idStr = fmt.Sprintf("%v", v)
					}

					options = append(options, map[string]interface{}{
						"id":   idStr,
						"name": name,
					})
				}
			}
		}

		return ctx.Respond(options, len(options))
	}

	return form.SelectField(id, title).
		Placeholder("Select product").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&listPrisyncProducts)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText(desc)
}

// GetBrandProp creates a dynamic select field for Prisync brands
func GetBrandProp(id string, title string, desc string, required bool, form *smartform.FormBuilder) *smartform.FieldBuilder {
	listPrisyncBrands := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		// Get API credentials from context
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		apiKey := authCtx.Extra["api-key"]
		apiToken := authCtx.Extra["api-token"]

		// Prisync API endpoint for fetching brands
		endpoint := "/api/v2/list/brand/startFrom/0"

		// Use the Prisync client to fetch brands
		response, err := PrisyncRequest(apiKey, apiToken, endpoint, http.MethodGet, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch Prisync brands: %v", err)
		}

		var options []map[string]interface{}

		// Parse the response based on Prisync API structure
		if responseMap, ok := response.(map[string]interface{}); ok {
			// Extract the results array (might be "results" or "brands" or at root level)
			var brands []interface{}

			// Try different possible response structures
			if results, ok := responseMap["results"].([]interface{}); ok {
				brands = results
			} else if brandsList, ok := responseMap["brands"].([]interface{}); ok {
				brands = brandsList
			} else if responseArray, ok := response.([]interface{}); ok {
				brands = responseArray
			}

			for _, item := range brands {
				brandMap, ok := item.(map[string]interface{})
				if !ok {
					continue
				}

				// Extract brand properties
				id, idOk := brandMap["id"]
				name, nameOk := brandMap["name"].(string)

				if !nameOk {
					name, nameOk = brandMap["brand_name"].(string)
				}

				if !idOk || !nameOk {
					continue
				}

				// Convert ID to string based on type
				var idStr string
				switch v := id.(type) {
				case float64:
					idStr = fmt.Sprintf("%.0f", v)
				case string:
					idStr = v
				case json.Number:
					idStr = string(v)
				default:
					idStr = fmt.Sprintf("%v", v)
				}

				options = append(options, map[string]interface{}{
					"id":   idStr,
					"name": name,
				})
			}
		}

		return ctx.Respond(options, len(options))
	}

	return form.SelectField(id, title).
		Placeholder("Select brand").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&listPrisyncBrands)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText(desc)
}

// GetCategoryProp creates a dynamic select field for Prisync categories
func GetCategoryProp(id string, title string, desc string, required bool, form *smartform.FormBuilder) *smartform.FieldBuilder {
	listPrisyncCategories := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		// Get API credentials from context
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		apiKey := authCtx.Extra["api-key"]
		apiToken := authCtx.Extra["api-token"]

		// Prisync API endpoint for fetching categories
		endpoint := "/api/v2/list/category/startFrom/0"

		// Use the Prisync client to fetch categories
		response, err := PrisyncRequest(apiKey, apiToken, endpoint, http.MethodGet, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch Prisync categories: %v", err)
		}

		// Process the response
		var options []map[string]interface{}

		// Parse the response based on Prisync API structure
		if responseMap, ok := response.(map[string]interface{}); ok {
			// Extract the results array (might be "results" or "categories" or at root level)
			var categories []interface{}

			// Try different possible response structures
			if results, ok := responseMap["results"].([]interface{}); ok {
				categories = results
			} else if categoriesList, ok := responseMap["categories"].([]interface{}); ok {
				categories = categoriesList
			} else if responseArray, ok := response.([]interface{}); ok {
				categories = responseArray
			}

			for _, item := range categories {
				categoryMap, ok := item.(map[string]interface{})
				if !ok {
					continue
				}

				// Extract category properties
				id, idOk := categoryMap["id"]
				name, nameOk := categoryMap["name"].(string)

				// Some APIs use "category_name" instead of "name"
				if !nameOk {
					name, nameOk = categoryMap["category_name"].(string)
				}

				if !idOk || !nameOk {
					continue
				}

				// Convert ID to string based on type
				var idStr string
				switch v := id.(type) {
				case float64:
					idStr = fmt.Sprintf("%.0f", v)
				case string:
					idStr = v
				case json.Number:
					idStr = string(v)
				default:
					idStr = fmt.Sprintf("%v", v)
				}

				options = append(options, map[string]interface{}{
					"id":   idStr,
					"name": name,
				})
			}
		}

		return ctx.Respond(options, len(options))
	}

	return form.SelectField(id, title).
		Placeholder("Select category").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&listPrisyncCategories)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText(desc)
}
