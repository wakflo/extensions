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

package freshdesk

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gookit/goutil/arrutil"
	fastshot "github.com/opus-domini/fast-shot"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var sharedAuth = autoform.NewCustomAuthField().
	SetFields(map[string]*sdkcore.AutoFormSchema{
		"domain": autoform.NewShortTextField().
			SetDisplayName("Freshdesk Domain").
			SetDescription("The domain name of the freshdesk account. eg. xyz.freshdesk.com, type in only 'xyz'").
			SetRequired(true).
			Build(),
		"api-key": autoform.NewShortTextField().SetDisplayName("Api Key").
			SetDescription("The api key used to authenticate freshdesk.").
			SetRequired(true).
			Build(),
	}).
	Build()

func NewFreshdeskAPIClient(baseURL, apiKey string) (*http.Client, string) {
	client := &http.Client{}

	// Encode the API key for basic authentication
	auth := base64.StdEncoding.EncodeToString([]byte(apiKey + ":X"))

	// Return the client and the base URL with the encoded authorization
	return client, auth
}

func GetTickets(baseURL, apiKey string) (interface{}, error) {
	client, auth := NewFreshdeskAPIClient(baseURL, apiKey)

	req, err := http.NewRequest(http.MethodGet, baseURL+"/api/v2/tickets", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+auth)

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
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return result, nil
}

func GetTicketQuery(baseURL, apiKey, date string) (interface{}, error) {
	client, auth := NewFreshdeskAPIClient(baseURL, apiKey)

	var urlSTR string
	if date != "" {
		query := url.QueryEscape(fmt.Sprintf(`"created_at:>'%s'"`, date))
		urlSTR = baseURL + "/api/v2/search/tickets?query=" + query
	} else {
		urlSTR = baseURL + "/api/v2/tickets"
	}
	req, err := http.NewRequest(http.MethodGet, urlSTR, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+auth)

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
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return result, nil
}

func CreateTicket(baseURL, apiKey string, ticketData map[string]interface{}) (map[string]interface{}, error) {
	client, auth := NewFreshdeskAPIClient(baseURL, apiKey)

	jsonData, err := json.Marshal(ticketData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/api/v2/tickets", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+auth)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var result map[string]interface{}
	if newError := json.Unmarshal(body, &result); newError != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", newError)
	}

	return result, nil
}

func UpdateTicket(baseURL, apiKey string, ticketID string, input TicketUpdate) error {
	client, auth := NewFreshdeskAPIClient(baseURL, apiKey)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v2/tickets/%s", baseURL, ticketID), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error fetching ticket: %s", resp.Status)
	}

	var existingTicket TicketUpdate
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	err = json.Unmarshal(body, &existingTicket)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	if input.Description != "" {
		existingTicket.Description = input.Description
	}
	if input.Subject != "" {
		existingTicket.Subject = input.Subject
	}
	if input.Priority != 0 {
		existingTicket.Priority = input.Priority
	}
	if input.Status != 0 {
		existingTicket.Status = input.Status
	}

	updateData, err := json.Marshal(existingTicket)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	updateReq, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/v2/tickets/%s", baseURL, ticketID), bytes.NewBuffer(updateData))
	if err != nil {
		return fmt.Errorf("failed to create update request: %v", err)
	}
	updateReq.Header.Set("Authorization", "Basic "+auth)
	updateReq.Header.Set("Content-Type", "application/json")

	updateResp, err := client.Do(updateReq)
	if err != nil {
		return fmt.Errorf("failed to send update request: %v", err)
	}
	defer updateResp.Body.Close()

	if updateResp.StatusCode != http.StatusOK {
		return fmt.Errorf("error updating ticket: %s", updateResp.Status)
	}

	return nil
}

func GetTicket(baseURL, apiKey, ticketID string) (interface{}, error) {
	client, auth := NewFreshdeskAPIClient(baseURL, apiKey)

	req, err := http.NewRequest(http.MethodGet, baseURL+"/api/v2/tickets/"+ticketID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+auth)

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
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return result, nil
}

func getTicketInput() *sdkcore.AutoFormSchema {
	getTickets := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
		baseAPI := "https://" + ctx.Auth.Extra["domain"] + ".freshdesk.com"
		qu := fastshot.NewClient(baseAPI).
			Auth().BasicAuth(ctx.Auth.Extra["api-key"], "X").
			Header().
			AddAccept("application/json").
			Build().GET("/api/v2/tickets")

		rsp, err := qu.Send()
		if err != nil {
			return nil, err
		}

		if rsp.Status().IsError() {
			return nil, errors.New(rsp.Status().Text())
		}
		defer rsp.Raw().Body.Close()

		bytes, err := io.ReadAll(rsp.Raw().Body) //nolint:bodyclose
		if err != nil {
			return nil, err
		}

		var tickets []Ticket
		err = json.Unmarshal(bytes, &tickets)
		if err != nil {
			return nil, err
		}

		return arrutil.Map[Ticket, map[string]any](tickets, func(input Ticket) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.ID,
				"name": input.Subject,
			}, true
		}), nil
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Ticket").
		SetDescription("Select ticket to update").
		SetDynamicOptions(&getTickets).
		SetRequired(true).Build()
}

var freshdeskPriorityType = []*sdkcore.AutoFormSchema{
	{Const: "4", Title: "Urgent"},
	{Const: "3", Title: "High"},
	{Const: "2", Title: "Normal"},
	{Const: "1", Title: "Low"},
}

var freshdeskStatusType = []*sdkcore.AutoFormSchema{
	{Const: "2", Title: "Open"},
	{Const: "3", Title: "Pending"},
	{Const: "4", Title: "Resolved"},
	{Const: "5", Title: "Closed"},
}