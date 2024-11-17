package typeform

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gookit/goutil/arrutil"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

const baseURL = "https://api.typeform.com/"

var (
	// #nosec
	tokenURL   = baseURL + "oauth/token"
	sharedAuth = autoform.NewOAuthField(baseURL+"oauth/authorize", &tokenURL, []string{
		"accounts:read",
		"forms:write",
		"forms:read",
		"responses:read",
		"responses:write",
	}).Build()
)

func getTypeformFormsInput(title string, desc string, required bool) *sdkcore.AutoFormSchema {
	getTypeformForms := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
		// Define the Typeform API URL for listing forms
		url := baseURL + "forms"

		// Create a new HTTP GET request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		// Set the required headers
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ctx.Auth.AccessToken))
		req.Header.Set("Content-Type", "application/json")

		// Create a new HTTP client and send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// Check if the response status indicates an error
		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("error: %s", resp.Status)
		}

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// Parse the response body into the struct
		var typeformResponse FormsResponse
		if err := json.Unmarshal(body, &typeformResponse); err != nil {
			return nil, err
		}

		// Extract forms from the response
		forms := typeformResponse.Items

		// Map the data into the expected format for AutoFormSchema
		return arrutil.Map[Form, map[string]any](forms, func(input Form) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.ID,
				"name": input.Title,
			}, true
		}), nil
	}

	// Return the AutoFormSchema using the dynamic form data
	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(desc).
		SetDynamicOptions(&getTypeformForms).
		SetRequired(required).
		Build()
}

func getFormResponses(accessToken, formID string) (map[string]interface{}, error) {
	// Define the API URL for retrieving form responses
	url := fmt.Sprintf(baseURL+"forms/%s/responses", formID)

	// Create a new HTTP GET request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Set required headers
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")

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
		return response, fmt.Errorf("error: %v", response)
	}

	return response, nil
}
