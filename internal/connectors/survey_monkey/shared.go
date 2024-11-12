package survey_monkey

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gookit/goutil/arrutil"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var (
	surveyMonkeyTokenURL = "https://api.surveymonkey.com/oauth/token"
	surveyMonkeyAuthURL  = "https://api.surveymonkey.com/oauth/authorize"

	// Setting up OAuthField for SurveyMonkey
	sharedAuth = autoform.NewOAuthField(
		surveyMonkeyAuthURL,
		&surveyMonkeyTokenURL,
		[]string{
			"contacts_write",
			"responses_read",
			"responses_read_detail",
			"webhooks_read",
			"webhooks_write",
		},
	).Build()
)

var baseURL = "https://api.surveymonkey.com/v3"

func createContactList(accessToken, listName string) (map[string]interface{}, error) {
	// Define the API URL for creating contact lists
	url := baseURL + "/contact_lists"

	// Set up the payload with the contact list name
	payload := map[string]interface{}{
		"name": listName,
	}

	// Marshal the payload into JSON format
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonPayload))
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
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return response, fmt.Errorf("error: %v", response)
	}

	return response, nil
}

func getSurveysInput(title string, desc string, required bool) *sdkcore.AutoFormSchema {
	getSurveys := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
		// Define the SurveyMonkey API URL for listing surveys
		url := baseURL + "/surveys"

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
		var surveysResponse SurveyMonkeySurveysResponse
		if err := json.Unmarshal(body, &surveysResponse); err != nil {
			return nil, err
		}

		// Extract surveys from the response
		surveys := surveysResponse.Data

		// Map the data into the expected format for AutoFormSchema
		return arrutil.Map[SurveyMonkeySurvey, map[string]any](surveys, func(input SurveyMonkeySurvey) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.ID,
				"name": input.Title,
			}, true
		}), nil
	}

	// Return the AutoFormSchema using the dynamic survey data
	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(desc).
		SetDynamicOptions(&getSurveys).
		SetRequired(required).
		Build()
}
