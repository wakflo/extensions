package actions

import (
	"fmt"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/campaignmonitor/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createCampaignFromTemplateActionProps struct {
	Name        string                   `json:"name"`
	Subject     string                   `json:"subject"`
	FromName    string                   `json:"fromName"`
	FromEmail   string                   `json:"fromEmail"`
	ReplyTo     string                   `json:"replyTo"`
	ListIDs     string                   `json:"listIds,omitempty"`
	SegmentIDs  []string                 `json:"segmentIds,omitempty"`
	TemplateID  string                   `json:"templateId"`
	Singlelines []map[string]string      `json:"singlelines,omitempty"`
	Multilines  []map[string]string      `json:"multilines,omitempty"`
	Images      []map[string]string      `json:"images,omitempty"`
	Repeaters   []map[string]interface{} `json:"repeaters,omitempty"`
}

type CreateCampaignFromTemplateAction struct{}

func (a *CreateCampaignFromTemplateAction) Name() string {
	return "Create Campaign From Template"
}

func (a *CreateCampaignFromTemplateAction) Description() string {
	return "Create a new draft campaign based on a template."
}

func (a *CreateCampaignFromTemplateAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateCampaignFromTemplateAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createCampaignTemplateDocs,
	}
}

func (a *CreateCampaignFromTemplateAction) Icon() *string {
	icon := "mdi:email-edit"
	return &icon
}

func (a *CreateCampaignFromTemplateAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"name": autoform.NewShortTextField().
			SetDisplayName("Campaign Name").
			SetDescription("The name of your campaign.").
			SetRequired(true).
			Build(),
		"subject": autoform.NewShortTextField().
			SetDisplayName("Email Subject").
			SetDescription("The subject line for your campaign.").
			SetRequired(true).
			Build(),
		"fromName": autoform.NewShortTextField().
			SetDisplayName("From Name").
			SetDescription("The name that will appear as the sender.").
			SetRequired(true).
			Build(),
		"fromEmail": autoform.NewShortTextField().
			SetDisplayName("From Email").
			SetDescription("The email address that will appear as the sender.").
			SetRequired(true).
			Build(),
		"replyTo": autoform.NewShortTextField().
			SetDisplayName("Reply-To Email").
			SetDescription("The email address that recipients can reply to.").
			SetRequired(true).
			Build(),
		"listIds": shared.GetCreateSendSubscriberListsInput(),
		"templateId": autoform.NewShortTextField().
			SetDisplayName("Template ID").
			SetDescription("The ID of the template to use.").
			SetRequired(true).
			Build(),
	}
}

func (a *CreateCampaignFromTemplateAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createCampaignFromTemplateActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	// Ensure at least one of ListIDs or SegmentIDs is provided
	if len(input.ListIDs) == 0 && len(input.SegmentIDs) == 0 {
		return nil, fmt.Errorf("at least one list ID or segment ID must be provided")
	}

	clientID := ctx.Auth.Extra["client-id"]
	if clientID == "" {
		return nil, fmt.Errorf("client ID is required")
	}

	payload := map[string]interface{}{
		"Name":       input.Name,
		"Subject":    input.Subject,
		"FromName":   input.FromName,
		"FromEmail":  input.FromEmail,
		"ReplyTo":    input.ReplyTo,
		"TemplateID": input.TemplateID,
	}

	if len(input.ListIDs) > 0 {
		payload["ListIDs"] = input.ListIDs
	}

	if len(input.SegmentIDs) > 0 {
		payload["SegmentIDs"] = input.SegmentIDs
	}

	templateContent := map[string]interface{}{}

	if len(input.Singlelines) > 0 {
		templateContent["Singlelines"] = input.Singlelines
	}

	if len(input.Multilines) > 0 {
		templateContent["Multilines"] = input.Multilines
	}

	if len(input.Images) > 0 {
		templateContent["Images"] = input.Images
	}

	if len(input.Repeaters) > 0 {
		templateContent["Repeaters"] = input.Repeaters
	}

	// Add template content to payload
	payload["TemplateContent"] = templateContent

	// Format the endpoint with the client ID
	endpoint := fmt.Sprintf("campaigns/%s/fromtemplate.json", clientID)

	// Make the API call to create the campaign
	result, err := shared.GetCampaignMonitorClient(
		ctx.Auth.Extra["api-key"],
		clientID,
		endpoint,
		http.MethodPost,
		payload)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *CreateCampaignFromTemplateAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateCampaignFromTemplateAction) SampleData() sdkcore.JSON {
	return "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1"
}

func (a *CreateCampaignFromTemplateAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateCampaignFromTemplateAction() sdk.Action {
	return &CreateCampaignFromTemplateAction{}
}
