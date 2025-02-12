package triggers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/wakflo/extensions/internal/integrations/zohoinventory/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type newPaymentTriggerProps struct {
	OrganizationID string `json:"organization_id"`
}

type NewPaymentTrigger struct{}

func (t *NewPaymentTrigger) Name() string {
	return "New Payment"
}

func (t *NewPaymentTrigger) Description() string {
	return "Triggered when a new payment is made, this integration allows you to automate workflows and processes in response to incoming payments, enabling seamless financial management and streamlined operations."
}

func (t *NewPaymentTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *NewPaymentTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &newPaymentDocs,
	}
}

func (t *NewPaymentTrigger) Icon() *string {
	return nil
}

func (t *NewPaymentTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"organization_id": shared.GetOrganizationsInput(),
	}
}

// Start initializes the newPaymentTrigger, required for event and webhook triggers in a lifecycle context.
func (t *NewPaymentTrigger) Start(ctx sdk.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the newPaymentTrigger, cleaning up resources and performing necessary teardown operations.
func (t *NewPaymentTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of newPaymentTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails. required for Pooling triggers
func (t *NewPaymentTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[newPaymentTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	var fromDate string
	lastRunTime := ctx.Metadata().LastRun

	if lastRunTime != nil {
		fromDate = lastRunTime.UTC().Format(time.RFC3339)
	}

	endpoint := shared.BaseURL + "/v1/customerpayments"

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %v", err)
	}

	q := u.Query()
	q.Set("organization_id", input.OrganizationID)
	q.Set("date", fromDate)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Zoho-oauthtoken "+ctx.Auth.Token.AccessToken)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

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

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return result, nil
}

func (t *NewPaymentTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *NewPaymentTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *NewPaymentTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func NewNewPaymentTrigger() sdk.Trigger {
	return &NewPaymentTrigger{}
}
