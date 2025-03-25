package shared

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var SharedAuth = autoform.NewAuth().NewCustomAuth().
	SetDescription("MailJet Authentication").
	SetLabel("Authentication details for connecting to MailJet").
	SetFields(map[string]*sdkcore.AutoFormSchema{
		"api_key": autoform.NewShortTextField().
			SetDisplayName("API Key").
			SetDescription("Your MailJet API Key").
			SetRequired(true).
			Build(),
		"secret_key": autoform.NewShortTextField().
			SetDisplayName("Secret Key").
			SetDescription("Your MailJet Secret Key").
			SetRequired(true).
			Build(),
	}).
	Build()

const (
	BaseURL = "https://api.mailjet.com"
	DataURL = "https://api.mailjet.com/v3.1/DATA"
)

// Client represents the MailJet API client
type Client struct {
	apiKey    string
	secretKey string
	client    *http.Client
}

// NewClient creates a new MailJet API client
func NewClient(apiKey, secretKey string) *Client {
	return &Client{
		apiKey:    apiKey,
		secretKey: secretKey,
		client:    &http.Client{},
	}
}

func (c *Client) Request(method, path string, payload interface{}, result interface{}) error {
	url := BaseURL + path

	var reqBody io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.apiKey, c.secretKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	errCode := 400
	if resp.StatusCode >= errCode {
		var errResp map[string]interface{}
		if err := json.Unmarshal(body, &errResp); err != nil {
			return fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
		}

		errMessage := fmt.Sprintf("API error: %d", resp.StatusCode)
		if errObj, ok := errResp["ErrorMessage"]; ok {
			errMessage = fmt.Sprintf("%s - %v", errMessage, errObj)
		}
		return errors.New(errMessage)
	}

	if result != nil {
		if err := json.Unmarshal(body, result); err != nil {
			return err
		}
	}

	return nil
}

func GetMailJetClient(apiKey, secretKey string) (*Client, error) {
	if apiKey == "" || secretKey == "" {
		return nil, errors.New("API key and Secret key are required")
	}
	return NewClient(apiKey, secretKey), nil
}
