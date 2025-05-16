// Copyright 2022-present Wakflo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shared

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gookit/goutil/arrutil"
	"github.com/juicycleff/smartform/v1"
	fastshot "github.com/opus-domini/fast-shot"
	"github.com/wakflo/go-sdk/v2"

	sdkcontext "github.com/wakflo/go-sdk/v2/context"

	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

const baseURL = "https://api.xero.com/api.xro/2.0"

var (
	// #nosec
	tokenURL = "https://identity.xero.com/connect/token"
	authURL  = "https://login.xero.com/identity/connect/authorize"
)

var (
	form = smartform.NewAuthForm("xero-auth", "Xero Oauth", smartform.AuthStrategyOAuth2)
	_    = form.OAuthField("oauth", "Xero Oauth").
		AuthorizationURL(authURL).
		TokenURL(tokenURL).
		Scopes([]string{"openid profile email accounting.transactions accounting.contacts accounting.attachments offline_access"}).
		Build()
)

var SharedAuth = form.Build()

// GetXeroNewClient sends a request to the Xero API using the provided access token.
func GetXeroNewClient(accessToken, endpoint, tenant string) (map[string]interface{}, error) {
	url := baseURL + endpoint

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Grant_type", "refresh_token")
	req.Header.Set("Xero-Tenant-Id", tenant)

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

	fmt.Printf("Response Status: %d\n", resp.StatusCode)
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

func GetTenantProps(id string, title string, desc string, required bool, form *smartform.FormBuilder) *smartform.FieldBuilder {
	getTenantID := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		tokenSource := ctx.Auth().Token
		if tokenSource == nil {
			return nil, errors.New("missing authentication token")
		}
		token := tokenSource.AccessToken

		client := fastshot.NewClient("https://api.xero.com").
			Auth().BearerToken(token).
			Header().
			AddAccept("application/json").
			Build()

		rsp, err := client.GET("/connections").Send()
		if err != nil {
			return nil, err
		}

		defer rsp.Body().Close()

		if rsp.Status().IsError() {
			return nil, errors.New(rsp.Status().Text())
		}

		defer rsp.Body().Close()
		byts, err := io.ReadAll(rsp.Body().Raw())
		if err != nil {
			return nil, err
		}

		var body TenantsResponse

		err = json.Unmarshal(byts, &body)
		if err != nil {
			return nil, err
		}

		sheet := body
		items := arrutil.Map[Tenant, map[string]any](sheet, func(input Tenant) (target map[string]any, find bool) {
			return map[string]any{
				"value": input.TenantID,
				"label": input.TenantName,
			}, true
		})

		return ctx.Respond(items, len(items))
	}

	return form.SelectField(id, title).
		Placeholder("Select Tenant").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getTenantID)).
				WithSearchSupport().
				End().
				GetDynamicSource(),
		).
		HelpText(desc)
}

func GetInvoiceProp(id string, title string, desc string, required bool, form *smartform.FormBuilder) *smartform.FieldBuilder {
	getInvoices := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		input := sdk.DynamicInputToType[struct {
			TenantID string `json:"tenant_id,omitempty"`
		}](ctx)

		endpoint := "https://api.xero.com/api.xro/2.0/Invoices"
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}
		tokenSource := ctx.Auth().Token
		if tokenSource == nil {
			return nil, errors.New("missing authentication token")
		}
		token := tokenSource.AccessToken

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Xero-Tenant-Id", input.TenantID)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to retrieve invoices, status code: %d, response: %s", resp.StatusCode, string(responseBody))
		}

		var result InvoicesResponse
		if err := json.Unmarshal(responseBody, &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %v", err)
		}

		invoice := result.Invoices
		items := arrutil.Map[Invoice, map[string]any](invoice, func(input Invoice) (target map[string]any, find bool) {
			return map[string]any{
				"value": input.InvoiceID,
				"label": input.InvoiceNumber,
			}, true
		})
		return ctx.Respond(items, len(items))
	}

	return form.SelectField(id, title).
		Placeholder("Select Invoice").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getInvoices)).
				WithSearchSupport().
				End().
				GetDynamicSource(),
		).
		HelpText(desc)
}

func SendInvoiceEmail(accessToken, endpoint, tenant string) error {
	url := baseURL + endpoint

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Grant_type", "refresh_token")

	req.Header.Set("Xero-Tenant-Id", tenant)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status code %d", resp.StatusCode)
	}

	fmt.Printf("Invoice email sent successfully\n")
	return nil
}

func CreateDraftInvoice(accessToken, tenant string, body map[string]interface{}) (sdkcore.JSON, error) {
	invoiceData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal invoice data: %v", err)
	}

	endpoint := "https://api.xero.com/api.xro/2.0/Invoices"
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(invoiceData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	req.Header.Set("Xero-Tenant-Id", tenant)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create draft invoice, status code: %d, response: %s", resp.StatusCode, string(responseBody))
	}

	var result map[string]interface{}
	if errs := json.Unmarshal(responseBody, &result); errs != nil {
		return nil, fmt.Errorf("failed to decode response: %v", errs)
	}

	return sdkcore.JSON(result), nil
}
