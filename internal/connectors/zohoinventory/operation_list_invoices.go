package zohoinventory

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type getInvoiceListOperationProps struct {
	OrganizationID string `json:"organization_id"`
}

type GetInvoiceListOperation struct {
	options *sdk.OperationInfo
}

func NewGetInvoiceListOperation() sdk.IOperation {
	return &GetInvoiceListOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Invoice List",
			Description: "Retrieve a list of Invoices",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"organization_id": autoform.NewShortTextField().
					SetDisplayName("Organization ID").
					SetDescription("The Zoho Inventory organization ID").
					SetRequired(true).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *GetInvoiceListOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing Zoho auth token")
	}

	input := sdk.InputToType[getInvoiceListOperationProps](ctx)

	url := "https://www.zohoapis.com/inventory/v1/invoices/?organization_id=" + input.OrganizationID

	req, err := http.NewRequest(http.MethodGet, url, nil)
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

func (c *GetInvoiceListOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetInvoiceListOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
