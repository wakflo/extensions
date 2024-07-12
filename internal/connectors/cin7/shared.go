package cin7

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var sharedAuth = autoform.NewCustomAuthField().
	SetFields(map[string]*sdkcore.AutoFormSchema{
		"account_id": autoform.NewShortTextField().SetDisplayName("Domain Name").
			SetDescription("Your Account ID").
			Build(),
		"key": autoform.NewShortTextField().SetDisplayName("Authentication Token").
			SetDescription("API Application Key").
			SetRequired(true).
			Build(),
	}).
	Build()

const baseURL = "https://inventory.dearsystems.com"

func fetchData(endpoint, accountID, applicationKey string, queryParams map[string]interface{}) (string, error) {
	params := url.Values{}
	for key, value := range queryParams {
		switch v := value.(type) {
		case string:
			params.Add(key, v)
		case int:
			params.Add(key, strconv.Itoa(v))
		case float64:
			params.Add(key, strconv.FormatFloat(v, 'f', -1, 64))
		case bool:
			params.Add(key, strconv.FormatBool(v))
		case time.Time:
			params.Add(key, v.Format(time.RFC3339))
		default:
			params.Add(key, fmt.Sprintf("%v", v))
		}
	}

	fullURL := fmt.Sprintf("%s%s?%s", baseURL, endpoint, params.Encode())

	client := &http.Client{}
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-auth-accountid", accountID)
	req.Header.Set("api-auth-applicationkey", applicationKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}
