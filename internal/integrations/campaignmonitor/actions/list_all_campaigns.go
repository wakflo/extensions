package actions

import (
	"fmt"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/campaignmonitor/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type ListAllCampaignsAction struct{}

func (a *ListAllCampaignsAction) Name() string {
	return "List All Campaigns"
}

func (a *ListAllCampaignsAction) Description() string {
	return "Retrieve all campaigns (sent, draft, and scheduled) from Campaign Monitor."
}

func (a *ListAllCampaignsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListAllCampaignsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listCampaignsDocs,
	}
}

func (a *ListAllCampaignsAction) Icon() *string {
	icon := "mdi:email-multiple-outline"
	return &icon
}

func (a *ListAllCampaignsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"page": autoform.NewNumberField().
			SetDisplayName("Page").
			SetDescription("The page number to retrieve (for pagination).").
			SetRequired(false).
			Build(),
	}
}

func (a *ListAllCampaignsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	clientID := ctx.Auth.Extra["client-id"]
	if clientID == "" {
		return nil, fmt.Errorf("client ID is required")
	}

	apiKey := ctx.Auth.Extra["api-key"]
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	// Combine results from all three endpoints
	campaigns := []interface{}{}

	// 1. Get sent campaigns
	sentEndpoint := fmt.Sprintf("clients/%s/campaigns.json", clientID)
	sentResult, err := shared.GetCampaignMonitorClient(
		apiKey,
		clientID,
		sentEndpoint,
		http.MethodGet,
		nil)

	if err == nil {
		if sentList, ok := sentResult.([]interface{}); ok {
			for _, sent := range sentList {
				if sentMap, ok := sent.(map[string]interface{}); ok {
					sentMap["Status"] = "Sent"
					campaigns = append(campaigns, sentMap)
				}
			}
		}
	}

	// 2. Get draft campaigns
	draftEndpoint := fmt.Sprintf("clients/%s/drafts.json", clientID)
	draftResult, err := shared.GetCampaignMonitorClient(
		apiKey,
		clientID,
		draftEndpoint,
		http.MethodGet,
		nil)

	if err == nil {
		if draftList, ok := draftResult.([]interface{}); ok {
			for _, draft := range draftList {
				if draftMap, ok := draft.(map[string]interface{}); ok {
					draftMap["Status"] = "Draft"
					campaigns = append(campaigns, draftMap)
				}
			}
		}
	}

	scheduledEndpoint := fmt.Sprintf("clients/%s/scheduled.json", clientID)
	scheduledResult, err := shared.GetCampaignMonitorClient(
		apiKey,
		clientID,
		scheduledEndpoint,
		http.MethodGet,
		nil)

	if err == nil {
		if scheduledList, ok := scheduledResult.([]interface{}); ok {
			for _, scheduled := range scheduledList {
				if scheduledMap, ok := scheduled.(map[string]interface{}); ok {
					scheduledMap["Status"] = "Scheduled"
					campaigns = append(campaigns, scheduledMap)
				}
			}
		}
	}

	return map[string]interface{}{
		"campaigns": campaigns,
		"count":     len(campaigns),
	}, nil
}

func (a *ListAllCampaignsAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListAllCampaignsAction) SampleData() sdkcore.JSON {
	return map[string]interface{}{
		"campaigns": []map[string]interface{}{
			{
				"CampaignID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
				"Name":       "Monthly Newsletter",
				"Subject":    "March 2025 Updates",
				"SentDate":   "2025-03-10 14:30",
				"Status":     "Sent",
				"FromName":   "Marketing Team",
				"FromEmail":  "marketing@example.com",
			},
			{
				"CampaignID": "b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2",
				"Name":       "Product Announcement",
				"Subject":    "Introducing Our New Product",
				"Status":     "Draft",
				"FromName":   "Product Team",
				"FromEmail":  "products@example.com",
			},
			{
				"CampaignID":    "c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3",
				"Name":          "Special Promotion",
				"Subject":       "25% Off This Weekend",
				"ScheduledDate": "2025-03-21 09:00",
				"Status":        "Scheduled",
				"FromName":      "Sales Team",
				"FromEmail":     "sales@example.com",
			},
		},
		"count": "3",
	}
}

func (a *ListAllCampaignsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListAllCampaignsAction() sdk.Action {
	return &ListAllCampaignsAction{}
}
