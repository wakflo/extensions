package actions

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/wakflo/extensions/internal/integrations/campaignmonitor/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type sendCampaignActionProps struct {
	CampaignID        string             `json:"campaignId"`
	ConfirmationEmail []shared.EmailItem `json:"confirmationEmail,omitempty"`
	SendDate          string             `json:"sendDate,omitempty"`
}

type SendCampaignAction struct{}

func (a *SendCampaignAction) Name() string {
	return "Send Campaign"
}

func (a *SendCampaignAction) Description() string {
	return "Schedule a draft campaign to be sent immediately or at a future date."
}

func (a *SendCampaignAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *SendCampaignAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &sendCampaignDocs,
	}
}

func (a *SendCampaignAction) Icon() *string {
	icon := "mdi:email-send"
	return &icon
}

func (a *SendCampaignAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"campaignId": autoform.NewShortTextField().
			SetDisplayName("Campaign ID").
			SetDescription("The ID of the draft campaign to send.").
			SetRequired(true).
			Build(),
		"confirmationEmail": autoform.NewArrayField().SetItems(
			autoform.NewShortTextField().
				SetDisplayName("Email").
				Build(),
		).
			SetDisplayName("Confirmation Emails").
			SetDescription("Email addresses to receive confirmation when the campaign is sent.").
			SetRequired(false).
			Build(),
		"sendDate": autoform.NewDateTimeField().
			SetDisplayName("Send Date").
			SetDescription("The date and time to send the campaign (format: YYYY-MM-DD HH:MM). Leave blank to send immediately.").
			SetRequired(false).
			Build(),
	}
}

func (a *SendCampaignAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[sendCampaignActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if input.CampaignID == "" {
		return nil, fmt.Errorf("campaign ID is required")
	}

	payload := map[string]interface{}{}

	if len(input.ConfirmationEmail) > 0 {
		emails := []string{}
		for _, item := range input.ConfirmationEmail {
			if item.Value != "" {
				emails = append(emails, item.Value)
			}
		}

		if len(emails) > 0 {
			payload["ConfirmationEmail"] = strings.Join(emails, ",")
		}
	}

	mark := 2
	if input.SendDate != "" {
		dateTimeParts := strings.Split(input.SendDate, "T")
		if len(dateTimeParts) > 1 {
			datePart := dateTimeParts[0]
			timePart := strings.Split(dateTimeParts[1], "+")[0]
			if len(strings.Split(timePart, ":")) > mark {
				timeComponents := strings.Split(timePart, ":")
				timePart = timeComponents[0] + ":" + timeComponents[1]
			}
			input.SendDate = datePart + " " + timePart
		}

		_, err := time.Parse("2006-01-02 15:04", input.SendDate)
		if err != nil {
			return nil, fmt.Errorf("invalid date format. Please use YYYY-MM-DD HH:MM format")
		}
		payload["SendDate"] = input.SendDate
	}

	endpoint := fmt.Sprintf("campaigns/%s/send.json", input.CampaignID)

	result, err := shared.GetCampaignMonitorClient(
		ctx.Auth.Extra["api-key"],
		ctx.Auth.Extra["client-id"],
		endpoint,
		http.MethodPost,
		payload)
	if err != nil {
		return nil, err
	}

	fmt.Println("Response:", result)

	return map[string]interface{}{
		"success":    true,
		"message":    "Campaign scheduled for sending",
		"CampaignID": input.CampaignID,
	}, nil
}

func (a *SendCampaignAction) Auth() *sdk.Auth {
	return nil
}

func (a *SendCampaignAction) SampleData() sdkcore.JSON {
	return map[string]interface{}{
		"success": true,
		"message": "Campaign scheduled for sending",
	}
}

func (a *SendCampaignAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewSendCampaignAction() sdk.Action {
	return &SendCampaignAction{}
}
