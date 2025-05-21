package shared

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"

	sdkcontext "github.com/wakflo/go-sdk/v2/context"

	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

var (
	form = smartform.NewAuthForm("mailjet-auth", "MailJet Authentication", smartform.AuthStrategyCustom)

	_ = form.TextField("api_key", "API Key").
		Placeholder("Enter API Key").
		Required(true).
		HelpText("Your MailJet API Key")

	_ = form.TextField("secret_key", "Secret Key").
		Placeholder("Enter Secret Key").
		Required(true).
		HelpText("Your MailJet Secret Key")

	SharedAuth = form.Build()
)

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

func GetContactProp(id string, title, desc string, required bool, form *smartform.FormBuilder) *smartform.FieldBuilder {
	getContacts := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		// Get the auth context
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		apiKey := authCtx.Extra["api_key"]
		secretKey := authCtx.Extra["secret_key"]

		client, err := GetMailJetClient(apiKey, secretKey)
		if err != nil {
			return nil, err
		}

		// Define a variable to hold the result
		var result map[string]interface{}

		// Make request to MailJet API with the correct function signature
		err = client.Request(http.MethodGet, "/v3/REST/contact", nil, &result)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch MailJet contacts: %v", err)
		}

		// Process the response
		var options []map[string]interface{}

		// Check if "Data" key exists and is an array
		contactsData, ok := result["Data"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected response format from MailJet API")
		}

		// Extract contact information from response
		for _, item := range contactsData {
			contactMap, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			// Extract contact properties
			id, idOk := contactMap["ID"]
			email, emailOk := contactMap["Email"].(string)

			if !idOk || !emailOk {
				continue
			}

			// Convert ID to string based on type
			var idStr string
			switch v := id.(type) {
			case float64:
				idStr = fmt.Sprintf("%.0f", v)
			case string:
				idStr = v
			case json.Number:
				idStr = string(v)
			default:
				idStr = fmt.Sprintf("%v", v)
			}

			options = append(options, map[string]interface{}{
				"id":   idStr,
				"name": email,
			})
		}

		return ctx.Respond(options, len(options))
	}

	return form.SelectField(id, title).
		Placeholder("Select contact").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getContacts)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText(desc)
}
