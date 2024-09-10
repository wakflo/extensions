package freshworkscrm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	sdk "github.com/wakflo/go-sdk/connector"

	"github.com/gookit/goutil/arrutil"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var sharedAuth = autoform.NewCustomAuthField().
	SetFields(map[string]*sdkcore.AutoFormSchema{
		"domain": autoform.NewShortTextField().
			SetDisplayName("Freshworks Domain").
			SetDescription("The domain name of the freshworks account. eg. xyz.freshworks.com, type in only 'xyz'").
			SetRequired(true).
			Build(),
		"api-key": autoform.NewShortTextField().SetDisplayName("Api Key").
			SetDescription("The api key used to authenticate freshworks crm.").
			SetRequired(true).
			Build(),
	}).
	Build()

func NewFreshWorksAPIClient(baseURL, apiKey string) *http.Client {
	return &http.Client{}
}

func CreateContact(baseURL, apiKey string, contactData map[string]interface{}) (interface{}, error) {
	client := NewFreshWorksAPIClient(baseURL, apiKey)

	jsonData, err := json.Marshal(contactData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/crm/sales/api/contacts", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token token="+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var result interface{}
	if newErr := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", newErr)
	}

	return result, nil
}

func getContactViewInput() *sdkcore.AutoFormSchema {
	getContactView := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
		baseAPI := "https://" + ctx.Auth.Extra["domain"] + ".myfreshworks.com"
		apiKey := ctx.Auth.Extra["api-key"]

		// Build the request
		req, err := http.NewRequest(http.MethodGet, baseAPI+"/crm/sales/api/contacts/filters", nil)
		if err != nil {
			return nil, err
		}

		// Set headers
		req.Header.Set("Authorization", "Token token="+apiKey)
		req.Header.Set("Accept", "application/json")

		// Send the request
		client := &http.Client{}
		rsp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer rsp.Body.Close()

		newBytes, err := io.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}

		var views FilterWrapper
		err = json.Unmarshal(newBytes, &views)
		if err != nil {
			return nil, err
		}

		view := views.Filters

		return arrutil.Map[ViewDetails, map[string]any](view, func(input ViewDetails) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.ID,
				"name": input.Name,
			}, true
		}), nil
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Contact View").
		SetDescription("Select to view list of contacts in a specific view").
		SetDynamicOptions(&getContactView).
		SetRequired(true).Build()
}

func getContactsInput() *sdkcore.AutoFormSchema {
	getContacts := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
		input := sdk.DynamicInputToType[struct {
			ContactViewID string `json:"contact_view_id"`
		}](ctx)
		baseAPI := "https://" + ctx.Auth.Extra["domain"] + ".myfreshworks.com"
		apiKey := ctx.Auth.Extra["api-key"]

		request := fmt.Sprintf("%s/crm/sales/api/contacts/view/%s", baseAPI, input.ContactViewID)

		req, err := http.NewRequest(http.MethodGet, request, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "Token token="+apiKey)
		req.Header.Set("Accept", "application/json")

		client := &http.Client{}
		rsp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer rsp.Body.Close()

		newBytes, err := io.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}

		var contacts ContactWrapper
		err = json.Unmarshal(newBytes, &contacts)
		if err != nil {
			return nil, err
		}

		contact := contacts.Contacts

		return arrutil.Map[Contact, map[string]any](contact, func(input Contact) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.ID,
				"name": input.DisplayName,
			}, true
		}), nil
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Contacts").
		SetDescription("Select a contact to update").
		SetDynamicOptions(&getContacts).
		SetRequired(true).Build()
}

func updateContact(baseURL, apiKey, contactID string, contactData map[string]interface{}) (interface{}, error) {
	client := NewFreshWorksAPIClient(baseURL, apiKey)

	jsonData, err := json.Marshal(contactData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, baseURL+"/crm/sales/api/contacts/"+contactID, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token token="+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var result interface{}
	if newErr := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", newErr)
	}

	return result, nil
}

func updateField(data map[string]interface{}, key, value string) {
	if value != "" {
		data[key] = value
	}
}
