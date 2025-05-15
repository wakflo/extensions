package shared

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

	"github.com/juicycleff/smartform/v1"
)

var (
	// #nosec
	tokenURL = baseURL + "/oauth2/token"
	authURL  = baseURL + "/oauth2/authorize"
)

var form = smartform.NewAuthForm("mailchimp-auth", "Mailchimp Oauth", smartform.AuthStrategyOAuth2)
var _ = form.OAuthField("oauth", "Mailchimp Oauth").
	AuthorizationURL(authURL).
	TokenURL(tokenURL).
	Scopes([]string{""}).
	Build()

var (
	SharedAuth = form.Build()
)

const baseURL = "https://login.mailchimp.com"

func GetMailChimpServerPrefix(accessToken string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/oauth2/metadata", nil)
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
		if errs := json.Unmarshal(body, &apiError); errs != nil {
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

func AddContactToList(accessToken, server, listID, email, firstName, status, lastName string) error {
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

var MailchimpStatusType = []*smartform.Option{
	{Value: "subscribed", Label: "Subscribed"},
	{Value: "unsubscribed", Label: "Unsubscribed"},
}

var MailchimpSubscriberStatus = []*smartform.Option{
	{Value: "subscribed", Label: "Subscribed"},
	{Value: "unsubscribed", Label: "Unsubscribed"},
	{Value: "cleaned", Label: "Cleaned"},
	{Value: "pending", Label: "Pending"},
	{Value: "transactional", Label: "Transactional"},
}

func FetchMailchimpLists(accessToken, server string) (interface{}, error) {
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

func UpdateSubscriberStatus(accessToken, server, listID, email, status string) error {
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

func ListRecentSubscribers(accessToken, server, listID, date string) (interface{}, error) {
	url := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/lists/%s/members?since_timestamp_opt=%s", server, listID, date)

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

	var data interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return data, nil
} // ListRecentSubscribers fetches and lists subscribers added or updated in the last 24 hours.
func ListRecentUnSubscribers(accessToken, server, listID, date string) (interface{}, error) {
	url := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/lists/%s/members?unsubscribed_since=%s", server, listID, date)

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

	var data interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return data, nil
}

func AddMemberNote(accessToken, server, listID, email, note string) error {
	subscriberHash := getSubscriberHash(email)

	url := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/lists/%s/members/%s/notes", server, listID, subscriberHash)

	payload := map[string]string{
		"note": note,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
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

// addSubscriberTag adds a tag to a subscriber in Mailchimp.
func ModifySubscriberTags(accessToken, server, listID, email string, tags []string, status string) error {
	subscriberHash := getSubscriberHash(email)

	url := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/lists/%s/members/%s/tags", server, listID, subscriberHash)

	tagEntries := make([]map[string]interface{}, 0, len(tags))

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

	if resp.StatusCode == http.StatusNoContent {
		fmt.Println("Action successful")
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	return nil
}

func ProcessTagNamesInput(input string) []string {
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
