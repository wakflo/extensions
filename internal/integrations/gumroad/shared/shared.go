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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var (
	// #nosec
	tokenURL   = "https://api.gumroad.com/oauth/token"
	SharedAuth = autoform.NewOAuthField("https://gumroad.com/oauth/authorize", &tokenURL, []string{
		"view_profile view_sales edit_products mark_sales_as_shipped refund_sales",
	}).Build()
)

const baseURL = "https://api.gumroad.com/v2"

func ListProducts(accessToken string, params url.Values) (map[string]interface{}, error) {

	// Define the API URL for fetching sales
	url := baseURL + "/products"

	// Append query parameters to URL if any exist
	if len(params) > 0 {
		url += "?" + params.Encode()
	}

	// Create a new HTTP GET request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Set required headers
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Initialize the HTTP client and execute the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response body
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	// Check for a successful response
	if res.StatusCode != http.StatusOK {
		return response, fmt.Errorf("error listing products, status code: %d, response: %v", res.StatusCode, response)
	}

	return response, nil

}

func GetProduct(accessToken string, productID string) (map[string]interface{}, error) {

	// Define the API URL for creating contact lists
	url := baseURL + "/products/" + productID

	// Append query parameters to URL if any exist
	// if len(params) > 0 {
	// 	url += "?" + params.Encode()
	// }

	// Create a new HTTP GET request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Set required headers
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Initialize the HTTP client and execute the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response body
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	// Check for a successful response
	if res.StatusCode != http.StatusOK {
		return response, fmt.Errorf("error listing products, status code: %d, response: %v", res.StatusCode, response)
	}

	return response, nil

}

func DisableProduct(accessToken string, productID string) (map[string]interface{}, error) {

	// Define the API URL for creating contact lists
	url := baseURL + "/products/" + productID + "/disable"

	// Append query parameters to URL if any exist
	// if len(params) > 0 {
	// 	url += "?" + params.Encode()
	// }

	// Create a new HTTP GET request
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return nil, err
	}

	// Set required headers
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Initialize the HTTP client and execute the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response body
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	// Check for a successful response
	if res.StatusCode != http.StatusOK {
		return response, fmt.Errorf("error listing products, status code: %d, response: %v", res.StatusCode, response)
	}

	return response, nil

}

func EnableProduct(accessToken string, productID string) (map[string]interface{}, error) {

	// Define the API URL for creating contact lists
	url := baseURL + "/products/" + productID + "/enable"

	// Append query parameters to URL if any exist
	// if len(params) > 0 {
	// 	url += "?" + params.Encode()
	// }

	// Create a new HTTP GET request
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return nil, err
	}

	// Set required headers
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Initialize the HTTP client and execute the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response body
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	// Check for a successful response
	if res.StatusCode != http.StatusOK {
		return response, fmt.Errorf("error listing products, status code: %d, response: %v", res.StatusCode, response)
	}

	return response, nil

}

func ListSales(accessToken string, params url.Values) (map[string]interface{}, error) {
	// Define the API URL for creating contact lists
	url := baseURL + "/sales"

	// Append query parameters to URL if any exist
	if len(params) > 0 {
		url += "?" + params.Encode()
	}

	fmt.Println("URL: ", url)
	fmt.Println("params: ", params)

	// Create a new HTTP GET request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Set required headers
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Initialize the HTTP client and execute the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response body
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	// Check for a successful response
	if res.StatusCode != http.StatusOK {
		return response, fmt.Errorf("error listing products, status code: %d, response: %v", res.StatusCode, response)
	}

	return response, nil

}

func GetSale(accessToken string, saleID string) (map[string]interface{}, error) {

	// Define the API URL for creating contact lists
	url := baseURL + "/sales/" + saleID

	// Create a new HTTP GET request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Set required headers
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Initialize the HTTP client and execute the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response body
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	// Check for a successful response
	if res.StatusCode != http.StatusOK {
		return response, fmt.Errorf("error listing products, status code: %d, response: %v", res.StatusCode, response)
	}

	return response, nil

}

func ListProductsInput(title string, desc string, required bool) *sdkcore.AutoFormSchema {
	listProducts := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		// Gumroad API endpoint for fetching products
		endpoint := "https://api.gumroad.com/v2/products"

		// Create a new HTTP GET request
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}

		// Set required headers according to Gumroad API documentation
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ctx.Auth.AccessToken)
		req.Header.Set("Accept", "application/json")

		// Create HTTP client and send request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		// Read response body
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}

		// Log the raw response for debugging
		fmt.Println("Raw API Response:", string(responseBody))

		// Parse the response
		var result ListProductsResponse
		if err := json.Unmarshal(responseBody, &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %v", err)
		}

		// Log the parsed result for debugging
		fmt.Printf("Parsed Result: %+v\n", result)

		// Map products to the format required by dynamic options
		var items []map[string]interface{}
		for _, product := range result.Products {
			items = append(items, map[string]interface{}{
				"id":   product.ID,
				"name": product.Name,
			})
		}

		// Log the items for debugging
		fmt.Printf("Items: %+v\n", items)

		return ctx.Respond(items, len(items))
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(desc).
		SetDynamicOptions(&listProducts).
		SetRequired(required).
		Build()
}

func FormatDateInput(dateStr string) (string, error) {
	// Remove the timezone ID part in square brackets
	cleanDateStr := strings.Split(dateStr, "[")[0]

	// Parse the date string
	t, err := time.Parse(time.RFC3339, cleanDateStr)
	if err != nil {
		return "", err
	}

	// Format to YYYY-MM-DD
	return t.Format("2006-01-02"), nil
}
