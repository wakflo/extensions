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

package triggers

import (
	"context"
	"log"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/cin7/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type newSalesTriggerProps struct {
	ID string `json:"id"`
}

type NewSalesTrigger struct{}

// Metadata returns metadata about the trigger
func (t *NewSalesTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "new_sales",
		DisplayName:   "New Sales",
		Description:   "Triggers workflow when a new sales is initiated",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: newSalesTriggerDocs,
		SampleOutput:  map[string]any{},
	}
}

// Auth returns the authentication requirements for the trigger
func (t *NewSalesTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

// GetType returns the type of trigger
func (t *NewSalesTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

// Props returns the schema for the trigger's input configuration
func (t *NewSalesTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("new_sales", "New Sales")

	form.TextField("id", "ID").
		Required(false).
		HelpText("Identifier for the trigger")

	schema := form.Build()

	return schema
}

// Start initializes the trigger
func (t *NewSalesTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Stop shuts down the trigger
func (t *NewSalesTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of the trigger
func (t *NewSalesTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	_, err := sdk.InputToTypeSafely[newSalesTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	endpoint := "/ExternalApi/SaleList"
	accountID := ctx.Auth().Extra["account_id"]
	applicationKey := ctx.Auth().Extra["key"]

	// Get last run time from metadata
	lastRun, err := ctx.GetMetadata("lastRun")
	if err != nil {
		return nil, err
	}

	var fromDate string
	if lastRunTime, ok := lastRun.(*time.Time); ok && lastRunTime != nil {
		fromDate = lastRunTime.UTC().Format(time.RFC3339)
	}

	queryParams := map[string]interface{}{
		"Page":         1,
		"CreatedSince": fromDate,
	}

	response, err := shared.FetchData(endpoint, accountID, applicationKey, queryParams)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	return response, nil
}

// Criteria defines the trigger criteria
func (t *NewSalesTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

// SampleData returns sample data for the trigger
func (t *NewSalesTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"sales": []map[string]any{
			{
				"id":           "123456",
				"name":         "Sample Sale",
				"customerName": "Sample Customer",
				"total":        100.00,
				"createdDate":  "2025-05-15T10:30:00Z",
			},
		},
	}
}

func NewNewSalesTrigger() sdk.Trigger {
	return &NewSalesTrigger{}
}
