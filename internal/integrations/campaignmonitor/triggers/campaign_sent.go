package triggers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/wakflo/extensions/internal/integrations/campaignmonitor/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type campaignSentTriggerProps struct {
	ClientID string `json:"clientId"`
}

type CampaignSentTrigger struct{}

func (t *CampaignSentTrigger) Name() string {
	return "Campaign Sent"
}

func (t *CampaignSentTrigger) Description() string {
	return "Trigger a workflow when a campaign is sent."
}

func (t *CampaignSentTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *CampaignSentTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &campaignSentDocs,
	}
}

func (t *CampaignSentTrigger) Icon() *string {
	icon := "mdi:email-check-outline"
	return &icon
}

func (t *CampaignSentTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"clientId": autoform.NewShortTextField().
			SetDisplayName("Client ID").
			SetDescription("The Client ID for which to monitor sent campaigns. If not provided, the Client ID from the authentication will be used.").
			SetRequired(false).
			Build(),
	}
}

// Start initializes the trigger
func (t *CampaignSentTrigger) Start(ctx sdk.LifecycleContext) error {
	return nil
}

// Stop shuts down the trigger
func (t *CampaignSentTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of CampaignSentTrigger
func (t *CampaignSentTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	// Get the input parameters
	input, err := sdk.InputToTypeSafely[campaignSentTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	// Get the last run time from metadata
	lastRunTime := ctx.Metadata().LastRun

	// Determine which client ID to use
	clientID := input.ClientID
	if clientID == "" {
		// Use the client ID from authentication if not provided in input
		clientID = ctx.Auth.Extra["client-id"]
		if clientID == "" {
			return nil, errors.New("client ID is required either as a parameter or in authentication")
		}
	}

	endpoint := fmt.Sprintf("clients/%s/campaigns.json", clientID)

	response, err := shared.GetCampaignMonitorClient(
		ctx.Auth.Extra["api-key"],
		clientID,
		endpoint,
		http.MethodGet,
		nil)
	if err != nil {
		return nil, err
	}

	campaignsList, ok := response.([]interface{})
	if !ok {
		return nil, errors.New("invalid response format: expected array of campaigns")
	}

	var newlySentCampaigns []interface{}

	for _, campaign := range campaignsList {
		campaignMap, ok := campaign.(map[string]interface{})
		if !ok {
			continue
		}

		sentDateStr, ok := campaignMap["SentDate"].(string)
		if !ok || sentDateStr == "" {
			continue
		}

		sentDate, err := time.Parse("2006-01-02 15:04", sentDateStr)
		if err != nil {
			sentDate, err = time.Parse(time.RFC3339, sentDateStr)
			if err != nil {
				continue
			}
		}

		if lastRunTime == nil || sentDate.After(*lastRunTime) {
			newlySentCampaigns = append(newlySentCampaigns, campaignMap)
		}
	}

	if len(newlySentCampaigns) == 0 {
		return map[string]interface{}{
			"newlySentCampaigns": []interface{}{},
			"count":              0,
			"clientId":           clientID,
		}, nil
	}

	return map[string]interface{}{
		"newlySentCampaigns": newlySentCampaigns,
		"count":              len(newlySentCampaigns),
		"clientId":           clientID,
	}, nil
}

func (t *CampaignSentTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *CampaignSentTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *CampaignSentTrigger) SampleData() sdkcore.JSON {
	return map[string]interface{}{
		"newlySentCampaigns": []interface{}{
			map[string]interface{}{
				"CampaignID":        "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
				"Name":              "Monthly Newsletter",
				"Subject":           "March 2025 Updates",
				"SentDate":          "2025-03-15 14:30",
				"FromName":          "Marketing Team",
				"FromEmail":         "marketing@example.com",
				"WebVersionURL":     "https://example.com/campaign/view",
				"WebVersionTextURL": "https://example.com/campaign/viewtext",
			},
			map[string]interface{}{
				"CampaignID":        "b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2",
				"Name":              "Product Announcement",
				"Subject":           "Introducing Our New Product",
				"SentDate":          "2025-03-16 09:45",
				"FromName":          "Product Team",
				"FromEmail":         "products@example.com",
				"WebVersionURL":     "https://example.com/campaign/view2",
				"WebVersionTextURL": "https://example.com/campaign/viewtext2",
			},
		},
		"count":    "2",
		"clientId": "c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3",
	}
}

func NewCampaignSentTrigger() sdk.Trigger {
	return &CampaignSentTrigger{}
}
