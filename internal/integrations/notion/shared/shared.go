package shared

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gookit/goutil/arrutil"
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

const BaseURL = "https://api.notion.com/v1"

// #nosec
var tokenURL = BaseURL + "oauth/token"

var form = smartform.NewAuthForm("notion-auth", "Notion Oauth", smartform.AuthStrategyOAuth2)

var _ = form.
	OAuthField("oauth", "Notion Oauth").
	AuthorizationURL("https://api.notion.com/v1/oauth/authorize").
	TokenURL("https://api.notion.com/v1/oauth/token").
	Scopes([]string{}).
	Build()

var SharedNotionAuth = form.Build()

func GetNotionPagesProp(title string, desc string, required bool, form *smartform.FormBuilder) *smartform.FieldBuilder {
	getPages := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		input := sdk.DynamicInputToType[struct {
			DatabaseID string `json:"databaseId"`
		}](ctx)

		log.Printf("[Notion] Fetching pages for database: %s", input.DatabaseID)

		// Validate that database ID is provided
		if input.DatabaseID == "" {
			log.Printf("[Notion] Error: Database ID is empty")
			return nil, errors.New("database ID is required to fetch pages")
		}

		client := &http.Client{}

		// Reference: https://developers.notion.com/reference/post-database-query
		// Endpoint: POST https://api.notion.com/v1/databases/{database_id}/query
		// The request body can contain filters, sorts, and pagination parameters
		// Even without filters, we need to send at least an empty JSON object
		url := fmt.Sprintf(BaseURL+"/databases/%s/query", input.DatabaseID)
		log.Printf("[Notion] Querying database at: %s", url)

		// Create request body - Notion API requires a JSON body even if empty
		// Reference: https://developers.notion.com/reference/post-database-query#body-parameters
		requestBody := map[string]interface{}{
			"page_size": 100, // Maximum results per page (default: 100, max: 100)
		}

		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			log.Printf("[Notion] Error marshaling request body: %v", err)
			return nil, err
		}

		log.Printf("[Notion] Request body: %s", string(jsonBody))

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			log.Printf("[Notion] Error creating HTTP request: %v", err)
			return nil, err
		}

		tokenSource := ctx.Auth().Token
		if tokenSource == nil {
			log.Printf("[Notion] Missing authentication token")
			return nil, errors.New("missing authentication token")
		}
		token := tokenSource.AccessToken

		// Set the required headers for the Notion API
		// Reference: https://developers.notion.com/reference/request-limits
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Notion-Version", "2022-06-28")
		req.Header.Set("Content-Type", "application/json")

		log.Printf("[Notion] Request headers set - Authorization: Bearer [REDACTED], Notion-Version: 2022-06-28")

		// Sending the request
		rsp, err := client.Do(req)
		if err != nil {
			log.Printf("[Notion] Error executing HTTP request: %v", err)
			return nil, err
		}

		defer rsp.Body.Close()

		log.Printf("[Notion] Response status code: %d", rsp.StatusCode)

		// Reading response body
		byts, err := io.ReadAll(rsp.Body)
		if err != nil {
			log.Printf("[Notion] Error reading response body: %v", err)
			return nil, err
		}

		log.Printf("[Notion] Response body: %s", string(byts))

		if rsp.StatusCode >= 400 {
			log.Printf("[Notion] API error - Status: %s, Body: %s", rsp.Status, string(byts))
			return nil, fmt.Errorf("notion API error %s: %s", rsp.Status, string(byts))
		}

		// Unmarshalling the response into a struct
		var body NotionQueryResponse
		err = json.Unmarshal(byts, &body)
		if err != nil {
			log.Printf("[Notion] Error unmarshaling response: %v", err)
			return nil, err
		}

		// Extracting pages from the response
		pages := body.Results
		log.Printf("[Notion] Found %d pages", len(pages))

		// Returning the mapped data in a format required by the AutoFormSchema
		items := arrutil.Map[NotionPage, map[string]any](pages, func(input NotionPage) (target map[string]any, find bool) {
			title := ""
			if len(input.Properties["Name"].Title) > 0 {
				title = input.Properties["Name"].Title[0].Text.Content
			}

			log.Printf("[Notion] Page: ID=%s, Title=%s", input.ID, title)

			return map[string]any{
				"id":   input.ID,
				"name": title,
				// "url":   input.URL,
			}, true
		})

		log.Printf("[Notion] Returning %d page options", len(items))
		return ctx.Respond(items, len(items))
	}

	return form.SelectField("page_id", title).
		Placeholder(desc).
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getPages)).
				WithSearchSupport().
				WithPagination(10).
				End().
				RefreshOn("databaseId").
				GetDynamicSource(),
		).
		HelpText("Select a page")
}

func GetNotionDatabasesProp(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getDatabases := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		// Reference: https://developers.notion.com/reference/post-search
		// Endpoint: POST https://api.notion.com/v1/search
		// This endpoint searches all pages and databases that have been shared with the integration
		url := BaseURL + "/search"

		log.Printf("[Notion] Fetching databases from: %s", url)

		// Create the request body with the filter to only get databases
		requestBody, err := json.Marshal(map[string]interface{}{
			"filter": map[string]interface{}{
				"property": "object",
				"value":    "database",
			},
		})
		if err != nil {
			log.Printf("[Notion] Error marshaling request body: %v", err)
			return nil, err
		}

		log.Printf("[Notion] Request body: %s", string(requestBody))

		// Create a new HTTP POST request
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
		if err != nil {
			log.Printf("[Notion] Error creating HTTP request: %v", err)
			return nil, err
		}

		tokenSource := ctx.Auth().Token
		if tokenSource == nil {
			log.Printf("[Notion] Missing authentication token")
			return nil, errors.New("missing authentication token")
		}
		token := tokenSource.AccessToken

		// Set the required headers
		// Reference: https://developers.notion.com/reference/versioning
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Notion-Version", "2022-06-28")     // Notion API version
		req.Header.Set("Content-Type", "application/json") // Content-Type for JSON

		log.Printf("[Notion] Request headers set - Authorization: Bearer [REDACTED], Notion-Version: 2022-06-28")

		// Create a new HTTP client and send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[Notion] Error executing HTTP request: %v", err)
			return nil, err
		}
		defer resp.Body.Close()

		log.Printf("[Notion] Response status code: %d", resp.StatusCode)

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[Notion] Error reading response body: %v", err)
			return nil, err
		}

		log.Printf("[Notion] Response body: %s", string(body))

		// Check if the response status indicates an error
		if resp.StatusCode >= 400 {
			log.Printf("[Notion] API error - Status: %s, Body: %s", resp.Status, string(body))
			return nil, fmt.Errorf("notion API error %s: %s", resp.Status, string(body))
		}

		// Parse the response body into the struct
		var bodyStruct NotionSearchResponse
		if err := json.Unmarshal(body, &bodyStruct); err != nil {
			log.Printf("[Notion] Error unmarshaling response: %v", err)
			return nil, err
		}

		// Extract databases from the response
		databases := bodyStruct.Results
		log.Printf("[Notion] Found %d databases", len(databases))

		// Map the data into the expected format for AutoFormSchema
		items := arrutil.Map[NotionDatabase, map[string]any](databases, func(input NotionDatabase) (target map[string]any, find bool) {
			title := ""
			if len(input.Title) > 0 && input.Title[0].Text.Content != "" {
				title = input.Title[0].Text.Content
			}

			log.Printf("[Notion] Database: ID=%s, Title=%s", input.ID, title)

			return map[string]any{
				"id":   input.ID,
				"name": title,
			}, true
		})

		log.Printf("[Notion] Returning %d database options", len(items))
		return ctx.Respond(items, len(items))
	}

	// Return the AutoFormSchema using the dynamic database data
	// Field name "databaseId" must match the JSON tag in action structs
	return form.SelectField("databaseId", "Database ID").
		Placeholder("Select a database").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getDatabases)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText("Select a database")
}

func CreateNotionPage(accessToken, parentPageID, title string, content string) (map[string]interface{}, error) {
	url := BaseURL + "/pages"

	// Create the payload with parent page ID and properties
	payload := map[string]interface{}{
		"parent": map[string]interface{}{
			"type":    "page_id",
			"page_id": parentPageID,
		},
		"properties": map[string]interface{}{
			"title": map[string]interface{}{
				"title": []map[string]interface{}{
					{
						"type": "text",
						"text": map[string]interface{}{
							"content": title,
						},
					},
				},
			},
		},
		"children": []map[string]interface{}{
			{
				"object": "block",
				"type":   "paragraph",
				"paragraph": map[string]interface{}{
					"rich_text": []map[string]interface{}{
						{
							"type": "text",
							"text": map[string]interface{}{
								"content": content,
							},
						},
					},
				},
			},
		},
	}

	// Marshal the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Create a new POST request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	// Set the required headers
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Notion-Version", "2022-06-28")

	// Create an HTTP client and send the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read and parse the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	// Check if the status code is not successful
	if res.StatusCode != http.StatusOK {
		return response, fmt.Errorf("error: %v", response)
	}

	return response, nil
}

func UpdateNotionPage(accessToken, pageID, title, content string) (map[string]interface{}, error) {
	// Reference: https://developers.notion.com/reference/patch-page
	// First, get the page to find the title property name
	page, err := GetNotionPage(accessToken, pageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get page: %w", err)
	}

	// Find the title property name from the page properties
	pageMap, ok := page.(map[string]interface{})
	if !ok {
		return nil, errors.New("page data is not in expected format")
	}

	properties, ok := pageMap["properties"].(map[string]interface{})
	if !ok {
		return nil, errors.New("page properties not found")
	}

	var titlePropertyName string
	for propName, propValue := range properties {
		if propMap, ok := propValue.(map[string]interface{}); ok {
			if propType, ok := propMap["type"].(string); ok && propType == "title" {
				titlePropertyName = propName
				break
			}
		}
	}

	if titlePropertyName == "" {
		return nil, errors.New("title property not found in page")
	}

	// Update the page title
	url := fmt.Sprintf(BaseURL+"/pages/%s", pageID)
	updatePayload := map[string]interface{}{
		"properties": map[string]interface{}{
			titlePropertyName: map[string]interface{}{
				"title": []map[string]interface{}{
					{
						"type": "text",
						"text": map[string]interface{}{
							"content": title,
						},
					},
				},
			},
		},
	}

	jsonPayload, err := json.Marshal(updatePayload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Notion-Version", "2022-06-28")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return response, fmt.Errorf("failed to update page (status %d): %v", res.StatusCode, response)
	}

	// Now append content as a new block to the page
	// Reference: https://developers.notion.com/reference/patch-block-children
	if content != "" {
		blockURL := fmt.Sprintf(BaseURL+"/blocks/%s/children", pageID)
		blockPayload := map[string]interface{}{
			"children": []map[string]interface{}{
				{
					"object": "block",
					"type":   "paragraph",
					"paragraph": map[string]interface{}{
						"rich_text": []map[string]interface{}{
							{
								"type": "text",
								"text": map[string]interface{}{
									"content": content,
								},
							},
						},
					},
				},
			},
		}

		blockJSON, err := json.Marshal(blockPayload)
		if err != nil {
			return response, nil // Return the page update even if content append fails
		}

		blockReq, err := http.NewRequest(http.MethodPatch, blockURL, bytes.NewBuffer(blockJSON))
		if err != nil {
			return response, nil
		}

		blockReq.Header.Add("Authorization", "Bearer "+accessToken)
		blockReq.Header.Add("Content-Type", "application/json")
		blockReq.Header.Add("Notion-Version", "2022-06-28")

		blockRes, err := client.Do(blockReq)
		if err != nil {
			return response, nil
		}
		defer blockRes.Body.Close()
	}

	return response, nil
}

func QueryNewPages(accessToken, databaseID string, lastChecked time.Time) ([]map[string]interface{}, error) {
	// Reference: https://developers.notion.com/reference/post-database-query
	// This function queries a database for pages created after a specific time
	log.Printf("[Notion] Querying new pages in database: %s since: %s", databaseID, lastChecked.Format(time.RFC3339))

	if databaseID == "" {
		log.Printf("[Notion] Error: Database ID is empty")
		return nil, errors.New("database ID is required")
	}

	url := fmt.Sprintf(BaseURL+"/databases/%s/query", databaseID)
	log.Printf("[Notion] POST request to: %s", url)

	payload := map[string]interface{}{
		"filter": map[string]interface{}{
			"property": "Created time",
			"date": map[string]interface{}{
				"after": lastChecked.Format(time.RFC3339),
			},
		},
		"sorts": []map[string]interface{}{
			{
				"property":  "Created time",
				"direction": "ascending",
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[Notion] Error marshaling request body: %v", err)
		return nil, err
	}

	log.Printf("[Notion] Request body: %s", string(jsonPayload))

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Printf("[Notion] Error creating HTTP request: %v", err)
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Notion-Version", "2022-06-28")

	log.Printf("[Notion] Request headers set - Authorization: Bearer [REDACTED], Notion-Version: 2022-06-28")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("[Notion] Error executing HTTP request: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	log.Printf("[Notion] Response status code: %d", res.StatusCode)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("[Notion] Error reading response body: %v", err)
		return nil, err
	}

	log.Printf("[Notion] Response body: %s", string(body))

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("[Notion] Error unmarshaling response: %v", err)
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("[Notion] API error - Status: %s, Body: %s", res.Status, string(body))
		return nil, fmt.Errorf("notion API error %s: %v", res.Status, response)
	}

	results, ok := response["results"].([]interface{})
	if !ok {
		log.Printf("[Notion] Error: unexpected response format - 'results' field missing or invalid")
		return nil, errors.New("unexpected response format")
	}

	log.Printf("[Notion] Found %d new pages", len(results))

	newPages := make([]map[string]interface{}, 0, len(results))
	for _, result := range results {
		page, ok := result.(map[string]interface{})
		if !ok {
			continue
		}
		newPages = append(newPages, map[string]interface{}{
			"id":               page["id"],
			"created_time":     page["created_time"],
			"last_edited_time": page["last_edited_time"],
			"url":              page["url"],
		})
	}

	log.Printf("[Notion] Returning %d new pages", len(newPages))
	return newPages, nil
}

func GetNotionPage(accessToken, pageID string) (sdkcore.JSON, error) {
	// Reference: https://developers.notion.com/reference/retrieve-a-page
	// Endpoint: GET https://api.notion.com/v1/pages/{page_id}
	log.Printf("[Notion] Retrieving page: %s", pageID)

	if pageID == "" {
		log.Printf("[Notion] Error: Page ID is empty")
		return nil, errors.New("page ID is required")
	}

	url := BaseURL + "/pages/" + pageID
	log.Printf("[Notion] GET request to: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[Notion] Error creating HTTP request: %v", err)
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Notion-Version", "2022-06-28")

	log.Printf("[Notion] Request headers set - Authorization: Bearer [REDACTED], Notion-Version: 2022-06-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[Notion] Error executing HTTP request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("[Notion] Response status code: %d", resp.StatusCode)

	// Read response body for better error messages
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[Notion] Error reading response body: %v", err)
		return nil, err
	}

	log.Printf("[Notion] Response body: %s", string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		log.Printf("[Notion] API error - Status: %s, Body: %s", resp.Status, string(bodyBytes))
		return nil, fmt.Errorf("notion API error %s: %s", resp.Status, string(bodyBytes))
	}

	var pageData sdkcore.JSON
	if err := json.Unmarshal(bodyBytes, &pageData); err != nil {
		log.Printf("[Notion] Error unmarshaling response: %v", err)
		return nil, err
	}

	log.Printf("[Notion] Successfully retrieved page: %s", pageID)
	return pageData, nil
}

// GetPageTitle Helper function to get the title from page properties
func GetPageTitle(properties map[string]interface{}) string {
	if titleProp, ok := properties["Name"].(map[string]interface{}); ok {
		if titleArray, ok := titleProp["title"].([]interface{}); ok && len(titleArray) > 0 {
			if titleText, ok := titleArray[0].(map[string]interface{}); ok {
				if plainText, ok := titleText["plain_text"].(string); ok {
					return plainText
				}
			}
		}
	}
	return "Untitled"
}
