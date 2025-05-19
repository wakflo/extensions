package shared

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
			DatabaseID string `json:"database"`
		}](ctx)

		client := &http.Client{}

		// Constructing the URL for querying the database
		url := fmt.Sprintf(BaseURL+"/databases/%s/query", input.DatabaseID)
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			return nil, err
		}

		tokenSource := ctx.Auth().Token
		if tokenSource == nil {
			return nil, errors.New("missing authentication token")
		}
		token := tokenSource.AccessToken

		// Set the required headers for the Notion API
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
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
		items := arrutil.Map[NotionPage, map[string]any](pages, func(input NotionPage) (target map[string]any, find bool) {
			title := ""
			if input.Properties["Name"].Title != nil && len(input.Properties["Name"].Title) > 0 {
				title = input.Properties["Name"].Title[0].Text.Content
			}

			return map[string]any{
				"id":   input.ID,
				"name": title,
				// "url":   input.URL,
			}, true
		})

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
				GetDynamicSource(),
		).
		HelpText("Select a page")
}

func GetNotionDatabasesProp(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getDatabases := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		// Define the Notion API URL
		url := BaseURL + "/v1/search"

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

		tokenSource := ctx.Auth().Token
		if tokenSource == nil {
			return nil, errors.New("missing authentication token")
		}
		token := tokenSource.AccessToken

		// Set the required headers
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
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
		items := arrutil.Map[NotionDatabase, map[string]any](databases, func(input NotionDatabase) (target map[string]any, find bool) {
			title := ""
			if len(input.Title) > 0 && input.Title[0].Text.Content != "" {
				title = input.Title[0].Text.Content
			}

			return map[string]any{
				"id":   input.ID,
				"name": title,
			}, true
		})
		return ctx.Respond(items, len(items))
	}

	// Return the AutoFormSchema using the dynamic database data
	return form.SelectField("database", "Database ID").
		Placeholder("Select a page").
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
		HelpText("Select a page")
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
	url := fmt.Sprintf(BaseURL+"/pages/%s", pageID)

	// Define the properties for title and content
	properties := map[string]interface{}{
		"Name": map[string]interface{}{
			"title": []map[string]interface{}{
				{
					"text": map[string]string{
						"content": title,
					},
				},
			},
		},
		"Content": map[string]interface{}{
			"rich_text": []map[string]interface{}{
				{
					"text": map[string]string{
						"content": content,
					},
				},
			},
		},
	}

	// Create the payload for updating the page
	payload := map[string]interface{}{
		"properties": properties,
	}

	// Convert the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	// Add headers
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Notion-Version", "2022-06-28") // Use the correct Notion API version

	// Send the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read the response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	// Check if the API returned an error
	if res.StatusCode != http.StatusOK {
		return response, fmt.Errorf("error: %v", response)
	}

	return response, nil
}

func QueryNewPages(accessToken, databaseID string, lastChecked time.Time) ([]map[string]interface{}, error) {
	url := fmt.Sprintf(BaseURL+"/databases/%s/query", databaseID)

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
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonPayload))
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
		return nil, fmt.Errorf("error: %v", response)
	}

	results, ok := response["results"].([]interface{})
	if !ok {
		return nil, errors.New("unexpected response format")
	}

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

	return newPages, nil
}

func GetNotionPage(accessToken, pageID string) (sdkcore.JSON, error) {
	url := BaseURL + "/pages/" + pageID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Notion-Version", "2022-06-28") // replace with the latest Notion API version

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to retrieve page from Notion")
	}

	var pageData sdkcore.JSON
	if err := json.NewDecoder(resp.Body).Decode(&pageData); err != nil {
		return nil, err
	}

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
