package xero

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

const baseURL = "https://api.xero.com/api.xro/2.0"

var (
	// #nosec
	tokenURL   = "https://identity.xero.com/connect/token"
	authURL    = "https://login.xero.com/identity/connect/authorize"
	sharedAuth = autoform.NewOAuthField(authURL, &tokenURL, []string{
		"openid profile email accounting.transactions accounting.contacts accounting.attachments offline_access",
	}).
		Build()
)

// getXeroNewClient sends a request to the Xero API using the provided access token.
func getXeroNewClient(accessToken, endpoint string) (map[string]interface{}, error) {
	url := baseURL + endpoint

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Grant_type", "refresh_token")

	tenantIDs, err := getTenantIDs(accessToken)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Xero-Tenant-Id", tenantIDs[0])

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

func sendInvoiceEmail(accessToken, endpoint string) error {
	url := baseURL + endpoint

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Grant_type", "refresh_token")
	tenantIDs, err := getTenantIDs(accessToken)
	if err != nil {
		return err
	}

	req.Header.Set("Xero-Tenant-Id", tenantIDs[0])

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

func createDraftInvoice(accessToken string, body map[string]interface{}) (sdk.JSON, error) {
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

	tenantIDs, err := getTenantIDs(accessToken)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Xero-Tenant-Id", tenantIDs[0])

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
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return sdk.JSON(result), nil
}

func getTenantIDs(accessToken string) ([]string, error) {
	url := "https://api.xero.com/connections"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Grant_type", "refresh_token")

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

	var tenants []map[string]interface{}
	if err := json.Unmarshal(body, &tenants); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	var tenantIDs []string
	for _, tenant := range tenants {
		if tenantID, ok := tenant["tenantId"].(string); ok {
			tenantIDs = append(tenantIDs, tenantID)
		}
	}

	return tenantIDs, nil
}

var xeroInvoiceStatus = []*sdkcore.AutoFormSchema{
	{Const: "DRAFT", Title: "Draft"},
	{Const: "SUBMITTED", Title: "Submitted"},
	{Const: "AUTHORISED", Title: "Authorised"},
	{Const: "DELETED", Title: "Delete"},
	{Const: "VOIDED", Title: "Voided"},
}
