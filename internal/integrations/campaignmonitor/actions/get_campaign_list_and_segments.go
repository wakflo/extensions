package actions

import (
	"fmt"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/campaignmonitor/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getCampaignListsAndSegmentsActionProps struct {
	CampaignID string `json:"campaignId"`
}

type GetCampaignListsAndSegmentsAction struct{}

func (a *GetCampaignListsAndSegmentsAction) Name() string {
	return "Get Campaign Lists and Segments"
}

func (a *GetCampaignListsAndSegmentsAction) Description() string {
	return "Retrieve the lists and segments a campaign was sent to."
}

func (a *GetCampaignListsAndSegmentsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetCampaignListsAndSegmentsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getCampaignListsAndSegmentsDocs,
	}
}

func (a *GetCampaignListsAndSegmentsAction) Icon() *string {
	icon := "mdi:account-group"
	return &icon
}

func (a *GetCampaignListsAndSegmentsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"campaignId": autoform.NewShortTextField().
			SetDisplayName("Campaign ID").
			SetDescription("The ID of the campaign to get lists and segments for.").
			SetRequired(true).
			Build(),
	}
}

func (a *GetCampaignListsAndSegmentsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getCampaignListsAndSegmentsActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if input.CampaignID == "" {
		return nil, fmt.Errorf("campaign ID is required")
	}

	// Format the endpoint with the campaign ID
	endpoint := fmt.Sprintf("campaigns/%s/listsandsegments.json", input.CampaignID)

	// Make the API call
	result, err := shared.GetCampaignMonitorClient(
		ctx.Auth.Extra["api-key"],
		ctx.Auth.Extra["client-id"],
		endpoint,
		http.MethodGet,
		nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *GetCampaignListsAndSegmentsAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetCampaignListsAndSegmentsAction) SampleData() sdkcore.JSON {
	return map[string]interface{}{
		"Lists": []map[string]interface{}{
			{
				"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
				"Name":   "My List 1",
			},
			{
				"ListID": "b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2",
				"Name":   "My List 2",
			},
		},
		"Segments": []map[string]interface{}{
			{
				"ListID":    "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
				"SegmentID": "c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3",
				"Title":     "My Segment 1",
			},
			{
				"ListID":    "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
				"SegmentID": "d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4",
				"Title":     "My Segment 2",
			},
		},
	}
}

func (a *GetCampaignListsAndSegmentsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetCampaignListsAndSegmentsAction() sdk.Action {
	return &GetCampaignListsAndSegmentsAction{}
}
