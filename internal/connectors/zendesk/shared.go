package zendesk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var sharedAuth = autoform.NewCustomAuthField().
	SetFields(map[string]*sdkcore.AutoFormSchema{
		"email": autoform.NewShortTextField().SetDisplayName("Agent Email (Required)").
			SetDescription("The email address you use to login to Zendesk").
			Build(),
		"api-token": autoform.NewShortTextField().SetDisplayName("Token (Required)").
			SetDescription("The API token you can generate in Zendesk.").
			Build(),
		"subdomain": autoform.NewShortTextField().SetDisplayName("Organization (e.g. wakflohelp) (Required)").
			SetDescription("The subdomain of your Zendesk instance.").
			Build(),
	}).
	Build()

func zendeskRequest(method, fullURL, email, apiToken string, request []byte) (map[string]interface{}, error) {
	req, err := http.NewRequest(method, fullURL, bytes.NewBuffer(request))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(apiToken+"%s/token", email)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

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
