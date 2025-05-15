package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"

	"github.com/gookit/goutil/arrutil"

	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

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

func RegisterContactViewProps(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getContactView := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		baseAPI := "https://" + authCtx.Extra["domain"] + ".myfreshworks.com"
		apiKey := authCtx.Extra["api-key"]

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

		items := arrutil.Map[ViewDetails, map[string]any](view, func(input ViewDetails) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.ID,
				"name": input.Name,
			}, true
		})

		return ctx.Respond(items, len(items))
	}

	return form.SelectField("contact_view_id", "Contact View").
		Placeholder("Select a contact view").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getContactView)).
				WithSearchSupport().
				End().
				GetDynamicSource(),
		).
		HelpText("Select to view list of contacts in a specific view")
}

func RegisterContactsProps(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getContacts := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		input := sdk.DynamicInputToType[struct {
			ContactViewID string `json:"contact_view_id"`
		}](ctx)

		baseAPI := "https://" + authCtx.Extra["domain"] + ".myfreshworks.com"
		apiKey := authCtx.Extra["api-key"]

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

		items := arrutil.Map[Contact, map[string]any](contact, func(input Contact) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.ID,
				"name": input.DisplayName,
			}, true
		})
		return ctx.Respond(items, len(items))
	}

	return form.SelectField("contact_id", "Contacts").
		Placeholder("Select a contact").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getContacts)).
				WithFieldReference("contact_view_id", "contactView").
				WithSearchSupport().
				End().
				RefreshOn("contactView").
				GetDynamicSource(),
		).
		HelpText("Select a contact to update")
}

func UpdateContact(baseURL, apiKey, contactID string, contactData map[string]interface{}) (interface{}, error) {
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

func ListContacts(baseURL, apiKey string, queryParams map[string]string) (interface{}, error) {
	client := NewFreshWorksAPIClient(baseURL, apiKey)

	endpoint := baseURL + "/crm/sales/api/contacts"
	if len(queryParams) > 0 {
		u, err := url.Parse(endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to parse URL: %v", err)
		}

		q := u.Query()
		for key, value := range queryParams {
			q.Add(key, value)
		}
		u.RawQuery = q.Encode()
		endpoint = u.String()
	}

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
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

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var result interface{}
	if errs := json.Unmarshal(body, &result); errs != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", errs)
	}

	return result, nil
}

func UpdateField(data map[string]interface{}, key, value string) {
	if value != "" {
		data[key] = value
	}
}
