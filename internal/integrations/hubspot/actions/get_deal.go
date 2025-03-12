package actions

import (
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/hubspot/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getDealActionProps struct {
	DealID string `json:"dealId"`
}

type GetDealAction struct{}

func (a *GetDealAction) Name() string {
	return "Get Deal"
}

func (a *GetDealAction) Description() string {
	return "Retrieve a specific HubSpot deal by ID"
}

func (a *GetDealAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetDealAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getDealDocs,
	}
}

func (a *GetDealAction) Icon() *string {
	return nil
}

func (a *GetDealAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"dealId": autoform.NewShortTextField().
			SetDisplayName("Deal ID").
			SetDescription("ID of the HubSpot deal to retrieve").
			SetRequired(true).Build(),
	}
}

func (a *GetDealAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getDealActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	reqURL := "/crm/v3/objects/deals/" + input.DealID

	resp, err := shared.HubspotClient(reqURL, ctx.Auth.AccessToken, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *GetDealAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetDealAction) SampleData() sdkcore.JSON {
	return map[string]interface{}{
		"id": "12345",
		"properties": map[string]interface{}{
			"dealname":            "Sample Deal",
			"amount":              "10000",
			"dealstage":           "presentationscheduled",
			"pipeline":            "default",
			"closedate":           "2023-12-31",
			"createdate":          "2023-01-15T09:30:00Z",
			"hs_lastmodifieddate": "2023-01-20T14:45:00Z",
		},
		"createdAt": "2023-01-15T09:30:00Z",
		"updatedAt": "2023-01-20T14:45:00Z",
		"archived":  false,
	}
}

func (a *GetDealAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetDealAction() sdk.Action {
	return &GetDealAction{}
}
