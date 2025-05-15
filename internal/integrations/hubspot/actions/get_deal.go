package actions

import (
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/hubspot/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getDealActionProps struct {
	DealID string `json:"dealId"`
}

type GetDealAction struct{}

// Metadata returns metadata about the action
func (a *GetDealAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_deal",
		DisplayName:   "Get Deal",
		Description:   "Retrieve a specific HubSpot deal by ID",
		Type:          core.ActionTypeAction,
		Documentation: getDealDocs,
		SampleOutput: map[string]interface{}{
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
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GetDealAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_deal", "Get Deal")

	form.TextField("dealId", "Deal ID").
		Required(true).
		HelpText("ID of the HubSpot deal to retrieve")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetDealAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GetDealAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getDealActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	reqURL := "/crm/v3/objects/deals/" + input.DealID

	resp, err := shared.HubspotClient(reqURL, authCtx.Token.AccessToken, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func NewGetDealAction() sdk.Action {
	return &GetDealAction{}
}
