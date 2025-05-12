package actions

import (
	"fmt"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/campaignmonitor/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getCampaignListsAndSegmentsActionProps struct {
	CampaignID string `json:"campaignId"`
}

type GetCampaignListsAndSegmentsAction struct{}

// Metadata returns metadata about the action
func (a *GetCampaignListsAndSegmentsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_campaign_lists_and_segments",
		DisplayName:   "Get Campaign Lists and Segments",
		Description:   "Retrieve the lists and segments a campaign was sent to.",
		Type:          core.ActionTypeAction,
		Documentation: getCampaignListsAndSegmentsDocs,
		Icon:          "mdi:account-group",
		SampleOutput: map[string]interface{}{
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
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GetCampaignListsAndSegmentsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_campaign_lists_and_segments", "Get Campaign Lists and Segments")

	form.TextField("campaignId", "Campaign ID").
		Placeholder("Enter campaign ID").
		Required(true).
		HelpText("The ID of the campaign to get lists and segments for.")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetCampaignListsAndSegmentsAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GetCampaignListsAndSegmentsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[getCampaignListsAndSegmentsActionProps](ctx)
	if err != nil {
		return nil, err
	}

	if input.CampaignID == "" {
		return nil, fmt.Errorf("campaign ID is required")
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	// Format the endpoint with the campaign ID
	endpoint := fmt.Sprintf("campaigns/%s/listsandsegments.json", input.CampaignID)

	// Make the API call
	result, err := shared.GetCampaignMonitorClient(
		authCtx.Extra["api-key"],
		authCtx.Extra["client-id"],
		endpoint,
		http.MethodGet,
		nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewGetCampaignListsAndSegmentsAction() sdk.Action {
	return &GetCampaignListsAndSegmentsAction{}
}
