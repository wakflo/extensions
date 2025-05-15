package triggers

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/stripe/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type newCustomerTriggerProps struct {
	Email          string     `json:"email,omitempty"`
	IncludeDeleted bool       `json:"includeDeleted"`
	CreatedTime    *time.Time `json:"createdTime"`
}

type NewCustomerTrigger struct{}

func (t *NewCustomerTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "new_customer",
		DisplayName:   "New Customer",
		Description:   "Triggers workflow when a new customer is created in Stripe",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: newCustomerDocs,
		SampleOutput: []map[string]any{
			{
				"id":          "cus_1234567890",
				"object":      "customer",
				"name":        "John Doe",
				"email":       "john.doe@example.com",
				"phone":       "+1234567890",
				"description": "New customer",
				"created":     1620000000,
				"address": map[string]string{
					"city":        "San Francisco",
					"country":     "US",
					"line1":       "123 Market St",
					"state":       "CA",
					"postal_code": "94102",
				},
			},
		},
	}
}

func (t *NewCustomerTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func (t *NewCustomerTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *NewCustomerTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("stripe-new-customer", "New Customer")

	form.TextField("email", "Email Filter").
		Placeholder("john@example.com").
		Required(false).
		HelpText("Only trigger for customers with this email address. Leave blank to trigger for all new customers.")

	form.CheckboxField("includeDeleted", "Include Deleted Customers").
		Placeholder("Include deleted customers").
		Required(false).
		HelpText("Whether to include deleted customers in the trigger.")

	schema := form.Build()

	return schema
}

// Start initializes the NewCustomerTrigger, required for event and webhook triggers in a lifecycle context.
func (t *NewCustomerTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the NewCustomerTrigger, cleaning up resources and performing necessary teardown operations.
func (t *NewCustomerTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of NewCustomerTrigger by processing the input context and returning a JSON response.
func (t *NewCustomerTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[newCustomerTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context and API key
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	apiKey := authCtx.Extra["api-key"]
	if apiKey == "" {
		return nil, errors.New("missing stripe secret api-key")
	}

	// Get the last run time from metadata
	var fromDate int64
	if input.CreatedTime == nil {
		lr, err := ctx.GetMetadata("lastRun")
		if err == nil && lr != nil {
			lastRunTime, ok := lr.(*time.Time)
			if ok && lastRunTime != nil {
				fromDate = lastRunTime.UTC().Unix()
			}
		}
	} else {
		fromDate = input.CreatedTime.UTC().Unix()
	}

	// Convert timestamp to string
	stringValue := strconv.FormatInt(fromDate, 10)

	// Set up query parameters
	params := url.Values{}
	if fromDate > 0 {
		params.Add("query", "created>='"+stringValue+"'")
	}
	if input.Email != "" {
		params.Add("email", input.Email)
	}
	if !input.IncludeDeleted {
		params.Add("deleted", "false")
	}

	reqURL := "/v1/customers"

	// Call the Stripe API
	resp, err := shared.StripeClient(apiKey, reqURL, http.MethodGet, nil, params)
	if err != nil {
		return nil, err
	}

	// Extract the data from the response
	nodes, ok := resp["data"].([]interface{})
	if !ok {
		return nil, errors.New("failed to extract data from response")
	}

	return nodes, nil
}

func (t *NewCustomerTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *NewCustomerTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":          "cus_1234567890",
		"object":      "customer",
		"name":        "John Doe",
		"email":       "john.doe@example.com",
		"phone":       "+1234567890",
		"description": "New customer",
		"created":     1620000000,
		"address": map[string]string{
			"city":        "San Francisco",
			"country":     "US",
			"line1":       "123 Market St",
			"state":       "CA",
			"postal_code": "94102",
		},
	}
}

func NewNewCustomerTrigger() sdk.Trigger {
	return &NewCustomerTrigger{}
}
