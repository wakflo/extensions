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

package hubspot

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	sdkcore "github.com/wakflo/go-sdk/core"

	"github.com/wakflo/go-sdk/autoform"
)

var (
	// #nosec
	tokenURL   = "https://api.hubapi.com/oauth/v1/token"
	sharedAuth = autoform.NewOAuthField("https://app.hubspot.com/oauth/authorize", &tokenURL, []string{
		"oauth" +
			" crm.lists.read " +
			"crm.lists.write " +
			"crm.objects.contacts.read" +
			" crm.objects.contacts.write " +
			"crm.objects.owners.read " +
			"crm.objects.companies.read " +
			"crm.objects.companies.write " +
			"crm.objects.deals.read" +
			" crm.objects.deals.write " +
			"crm.objects.line_items.read " +
			"crm.schemas.line_items.read " +
			"crm.schemas.companies.read " +
			"crm.schemas.contacts.read " +
			"crm.schemas.deals.read tickets",
	}).SetRequired(true).Build()
)

var baseAPI = "https://api.hubapi.com"

func hubspotClient(reqURL, accessToken, method string, request []byte) (interface{}, error) {
	fullURL := baseAPI + reqURL
	req, err := http.NewRequest(method, fullURL, bytes.NewBuffer(request))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var hubspotResponse interface{}
	err = json.Unmarshal(body, &hubspotResponse)
	if err != nil {
		return nil, err
	}

	return hubspotResponse, nil
}

//	func getListsInput() *sdkcore.AutoFormSchema {
//		getLists := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
//			client := fastshot.NewClient(baseAPI).
//				Auth().BearerToken(ctx.Auth.AccessToken).
//				Header().
//				AddAccept("application/json").
//				Build()
//
//			rsp, err := client.GET("/crm/v3/lists/").Send()
//			if err != nil {
//				return nil, err
//			}
//
//			if rsp.Status().IsError() {
//				return nil, errors.New(rsp.Status().Text())
//			}
//
//			bytes, err := io.ReadAll(rsp.Raw().Body) //nolint:bodyclose
//			if err != nil {
//				return nil, err
//			}
//
//			var hubspotLists ListResponse
//			err = json.Unmarshal(bytes, &hubspotLists)
//			if err != nil {
//				return nil, err
//			}
//
//			list := hubspotLists.Lists
//			return arrutil.Map[List, map[string]any](list, func(input List) (target map[string]any, find bool) {
//				return map[string]any{
//					"id":   input.ListID,
//					"name": input.Name,
//				}, true
//			}), nil
//		}
//
//		return autoform.NewDynamicField(sdkcore.String).
//			SetDisplayName("List").
//			SetDescription("Select list").
//			SetDependsOn([]string{"connection"}).
//			SetDynamicOptions(&getLists).
//			SetRequired(true).Build()
//	}
//
//	 func getPipelineInput() *sdkcore.AutoFormSchema {
//		getPipelines := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
//			client := fastshot.NewClient(baseAPI).
//				Auth().BearerToken(ctx.Auth.AccessToken).
//				Header().
//				AddAccept("application/json").
//				Build()
//
//			rsp, err := client.GET("/crm/v3/pipelines/deals").Send()
//			if err != nil {
//				return nil, err
//			}
//
//			if rsp.Status().IsError() {
//				return nil, errors.New(rsp.Status().Text())
//			}
//
//			bytes, err := io.ReadAll(rsp.Raw().Body) //nolint:bodyclose
//			if err != nil {
//				return nil, err
//			}
//
//			var hubspotPipelines DealPipelineResponse
//			err = json.Unmarshal(bytes, &hubspotPipelines)
//			if err != nil {
//				return nil, err
//			}
//
//			pipeline := hubspotPipelines.Results
//			return arrutil.Map[PipelineResult, map[string]any](pipeline, func(input PipelineResult) (target map[string]any, find bool) {
//				return map[string]any{
//					"id":   input.ID,
//					"name": input.Label,
//				}, true
//			}), nil
//		}
//
//		return autoform.NewDynamicField(sdkcore.String).
//			SetDisplayName("Deal Pipeline").
//			SetDescription("Deal Pipeline").
//			SetDependsOn([]string{"connection"}).
//			SetDynamicOptions(&getPipelines).
//			SetRequired(true).Build()
//	}

var hubspotPriority = []*sdkcore.AutoFormSchema{
	{Const: "HIGH", Title: "High"},
	{Const: "MEDIUM", Title: "Medium"},
	{Const: "LOW", Title: "Low"},
}
