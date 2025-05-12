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

type ListAllCampaignsAction struct{}

// Metadata returns metadata about the action
func (a *ListAllCampaignsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_all_campaigns",
		DisplayName:   "List All Campaigns",
		Description:   "Retrieve all campaigns (sent, draft, and scheduled) from Campaign Monitor.",
		Type:          core.ActionTypeAction,
		Documentation: listCampaignsDocs,
		Icon:          "mdi:email-multiple-outline",
		SampleOutput: map[string]interface{}{
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
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ListAllCampaignsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_all_campaigns", "List All Campaigns")

	form.NumberField("page", "Page").
		Placeholder("Enter page number").
		Required(false).
		HelpText("The page number to retrieve (for pagination).")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ListAllCampaignsAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ListAllCampaignsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	clientID := authCtx.Extra["client-id"]
	if clientID == "" {
		return nil, fmt.Errorf("client ID is required")
	}

	apiKey := authCtx.Extra["api-key"]
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

func NewListAllCampaignsAction() sdk.Action {
	return &ListAllCampaignsAction{}
}
