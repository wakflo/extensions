package mailchimp

import (
	"bytes"
	// #nosec
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	// #nosec

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var (
	// #nosec
	tokenURL   = "https://login.mailchimp.com/oauth2/token"
	authURL    = "https://login.mailchimp.com/oauth2/authorize"
	sharedAuth = autoform.NewOAuthField(authURL, &tokenURL, []string{
		"",
	}).
		Build()
)

func getMailChimpServerPrefix(accessToken string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://login.mailchimp.com/oauth2/metadata", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "OAuth "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Printf("Metadata response body: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		var apiError map[string]interface{}
		if err := json.Unmarshal(body, &apiError); err != nil {
			return "", fmt.Errorf("failed to unmarshal error response: %w", err)
		}
		return "", fmt.Errorf("API error: %v", apiError)
	}

	var metadata struct {
		DC string `json:"dc"`
	}
	if err := json.Unmarshal(body, &metadata); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if metadata.DC == "" {
		return "", errors.New("data center (dc) not found in the response")
	}

	return metadata.DC, nil
}

//  func getListInput() *sdkcore.AutoFormSchema {
//	getLists := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
//		dc, err := getMailChimpServerPrefix(ctx.Auth.AccessToken)
//		if err != nil {
//			log.Fatalf("Error getting MailChimp server prefix: %v", err)
//		}
//
//		url := fmt.Sprintf("https://%s.api.mailchimp.com/3.0", dc)
//
//		client := fastshot.NewClient(url).
//			Auth().BearerToken(ctx.Auth.AccessToken).
//			Header().
//			AddAccept("application/json").
//			Build()
//
//		rsp, err := client.GET("/lists").Send()
//		if err != nil {
//			return nil, err
//		}
//
//		if rsp.Status().IsError() {
//			return nil, errors.New(rsp.Status().Text())
//		}
//
//		body, err := io.ReadAll(rsp.Raw().Body) //nolint:bodyclose
//		if err != nil {
//			return nil, err
//		}
//
//		var lists ListResponse
//		err = json.Unmarshal(body, &lists)
//		if err != nil {
//			return nil, err
//		}
//
//		list := lists.Lists
//		return list, nil
//	}
//
//	return autoform.NewDynamicField(sdkcore.String).
//		SetDisplayName("Lists").
//		SetDescription("Select a list").
//		SetDependsOn([]string{"connection"}).
//		SetDynamicOptions(&getLists).
//		SetRequired(true).Build()
//  }

func addContactToList(accessToken, server, listID, email, firstName, status, lastName string) error {
	url := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/lists/%s/members", server, listID)

	payload := map[string]interface{}{
		"email_address": email,
		"status":        status,
		"merge_fields": map[string]string{
			"FNAME": firstName,
			"LNAME": lastName,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	return nil
}

var mailchimpStatusType = []*sdkcore.AutoFormSchema{
	{Const: "subscribed", Title: "Subscribed"},
	{Const: "unsubscribed", Title: "Unsubscribed"},
}

var mailchimpSubscriberStatus = []*sdkcore.AutoFormSchema{
	{Const: "subscribed", Title: "Subscribed"},
	{Const: "unsubscribed", Title: "Unsubscribed"},
	{Const: "cleaned", Title: "Cleaned"},
	{Const: "pending", Title: "Pending"},
	{Const: "transactional", Title: "Transactional"},
}

// FetchMailchimpLists retrieves the lists (audiences) from Mailchimp
func fetchMailchimpLists(accessToken, server string) (interface{}, error) {
	url := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/lists", server)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}

// updateSubscriberStatus updates the status of a subscriber in Mailchimp.
func updateSubscriberStatus(accessToken, server, listID, email, status string) error {
	subscriberHash := getSubscriberHash(email)

	url := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/lists/%s/members/%s", server, listID, subscriberHash)

	payload := map[string]interface{}{
		"status": status,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	return nil
}

// listRecentSubscribers fetches and lists subscribers added or updated in the last 24 hours.
func listRecentSubscribers(accessToken, server, listID, date string) (interface{}, error) {
	url := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/lists/%s/members?since_timestamp_opt=%s", server, listID, date)

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the JSON response
	var data interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return data, nil
} // listRecentSubscribers fetches and lists subscribers added or updated in the last 24 hours.
func listRecentUnSubscribers(accessToken, server, listID, date string) (interface{}, error) {
	url := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/lists/%s/members?unsubscribed_since=%s", server, listID, date)

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the JSON response
	var data interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return data, nil
}

func addMemberNote(accessToken, server, listID, email, note string) error {
	subscriberHash := getSubscriberHash(email)

	url := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/lists/%s/members/%s/notes", server, listID, subscriberHash)

	payload := map[string]string{
		"note": note,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	return nil
}

// addSubscriberTag adds a tag to a subscriber in Mailchimp.
func modifySubscriberTags(accessToken, server, listID, email string, tags []string, status string) error {
	// Calculate the subscriber hash from the email address
	subscriberHash := getSubscriberHash(email)

	url := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/lists/%s/members/%s/tags", server, listID, subscriberHash)

	tagEntries := make([]map[string]interface{}, 0, len(tags))

	// Prepare the request payload
	for _, tag := range tags {
		tagEntries = append(tagEntries, map[string]interface{}{
			"name":   tag,
			"status": status,
		})
	}

	payload := map[string]interface{}{
		"tags": tagEntries,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		// Consider the operation successful
		fmt.Println("Operation successful")
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	return nil
}

func processTagNamesInput(input string) []string {
	tagNames := strings.Split(input, ",")
	for i, tag := range tagNames {
		tagNames[i] = strings.TrimSpace(tag)
	}
	return tagNames
}

func getSubscriberHash(email string) string {
	emailLower := strings.ToLower(email)
	// #nosec
	hash := md5.Sum([]byte(emailLower))
	// #nosec
	return hex.EncodeToString(hash[:])
	// #nosec
}
