package zohoinventory

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
)

var (
	// #nosec
	tokenURL   = "https://accounts.zoho.com/oauth/v2/token"
	authURL    = "https://accounts.zoho.com/oauth/v2/auth"
	sharedAuth = autoform.NewOAuthField(authURL, &tokenURL, []string{
		"ZohoInventory.FullAccess.all",
	}).
		Build()
)

func getZohoClient(accessToken, url string) (map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Zoho-oauthtoken "+accessToken)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return result, nil
}
