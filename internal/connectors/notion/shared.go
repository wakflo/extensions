package notion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gookit/goutil/arrutil"
	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
	"io"
	"net/http"
)

var (
	// #nosec
	tokenURL   = "https://api.notion.com/v1/oauth/token"
	sharedAuth = autoform.NewOAuthField("https://api.notion.com/v1/oauth/authorize", &tokenURL, []string{}).Build()
)

const baseURL = "https://api.notion.com/v1"

func getNotionPagesInput(title string, desc string, required bool, databaseID string) *sdkcore.AutoFormSchema {
	getPages := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {

		input := sdk.DynamicInputToType[struct {
			DatabaseID string `json:"database"`
		}](ctx)

		client := &http.Client{}

		// Constructing the URL for querying the database
		url := fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", input.DatabaseID)
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			return nil, err
		}

		// Set the required headers for the Notion API
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ctx.Auth.AccessToken))
		req.Header.Set("Notion-Version", "2022-06-28")
		req.Header.Set("Content-Type", "application/json")

		// Sending the request
		rsp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		defer rsp.Body.Close()

		if rsp.StatusCode >= 400 {
			return nil, fmt.Errorf("error: %s", rsp.Status)
		}

		// Reading response body
		byts, err := io.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}

		// Unmarshalling the response into a struct
		var body NotionQueryResponse
		err = json.Unmarshal(byts, &body)
		if err != nil {
			return nil, err
		}

		// Extracting pages from the response
		pages := body.Results

		// Returning the mapped data in a format required by the AutoFormSchema
		return arrutil.Map[NotionPage, map[string]any](pages, func(input NotionPage) (target map[string]any, find bool) {
			title := ""
			if input.Properties["Name"].Title != nil && len(input.Properties["Name"].Title) > 0 {
				title = input.Properties["Name"].Title[0].Text.Content
			}

			return map[string]any{
				"id":   input.ID,
				"name": title,
				//"url":   input.URL,
			}, true
		}), nil
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(desc).
		SetDynamicOptions(&getPages).
		SetRequired(required).
		Build()
}

func getNotionDatabasesInput(title string, desc string, required bool) *sdkcore.AutoFormSchema {
	getDatabases := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
		// Define the Notion API URL
		url := "https://api.notion.com/v1/search"

		//input := sdk.DynamicInputToType[struct {
		//	User string `json:"user"`
		//}](ctx)

		// Create the request body with the filter to only get databases
		requestBody, err := json.Marshal(map[string]interface{}{
			"filter": map[string]interface{}{
				"property": "object",
				"value":    "database",
			},
		})
		if err != nil {
			return nil, err
		}

		// Create a new HTTP POST request
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
		if err != nil {
			return nil, err
		}

		// Set the required headers
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ctx.Auth.AccessToken))
		req.Header.Set("Notion-Version", "2022-06-28")     // Notion API version
		req.Header.Set("Content-Type", "application/json") // Content-Type for JSON

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
		var bodyStruct NotionSearchResponse
		if err := json.Unmarshal(body, &bodyStruct); err != nil {
			return nil, err
		}

		// Extract databases from the response
		databases := bodyStruct.Results

		// Map the data into the expected format for AutoFormSchema
		return arrutil.Map[NotionDatabase, map[string]any](databases, func(input NotionDatabase) (target map[string]any, find bool) {
			title := ""
			if len(input.Title) > 0 && input.Title[0].Text.Content != "" {
				title = input.Title[0].Text.Content
			}

			return map[string]any{
				"id":    input.ID,
				"title": title,
			}, true
		}), nil
	}

	// Return the AutoFormSchema using the dynamic database data
	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(desc).
		SetDynamicOptions(&getDatabases).
		SetRequired(required).
		Build()
}
