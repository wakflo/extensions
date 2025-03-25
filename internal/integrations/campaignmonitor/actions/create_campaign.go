package actions

import (
	"fmt"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/campaignmonitor/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createCampaignActionProps struct {
	Name       string   `json:"name"`
	Subject    string   `json:"subject"`
	FromName   string   `json:"fromName"`
	FromEmail  string   `json:"fromEmail"`
	ReplyTo    string   `json:"replyTo"`
	HtmlURL    string   `json:"htmlUrl"`
	TextURL    string   `json:"textUrl,omitempty"`
	ListID     string   `json:"listId,omitempty"`
	SegmentIDs []string `json:"segmentIds,omitempty"`
	InlineCSS  *bool    `json:"inlineCss,omitempty"`
}

type CreateCampaignAction struct{}

func (a *CreateCampaignAction) Name() string {
	return "Create Campaign"
}

func (a *CreateCampaignAction) Description() string {
	return "Create a new draft campaign ready to be tested or sent."
}

func (a *CreateCampaignAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateCampaignAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createCampaignDocs,
	}
}

func (a *CreateCampaignAction) Icon() *string {
	icon := "mdi:email-newsletter"
	return &icon
}

func (a *CreateCampaignAction) Properties() map[string]*sdkcore.AutoFormSchema {
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
		"htmlUrl": autoform.NewShortTextField().
			SetDisplayName("HTML URL").
			SetDescription("The URL of the HTML content for your campaign.").
			SetRequired(true).
			Build(),
		"textUrl": autoform.NewShortTextField().
			SetDisplayName("Text URL").
			SetDescription("The URL of the text content for your campaign (optional).").
			SetRequired(false).
			Build(),
		"listId": shared.GetCreateSendSubscriberListsInput(),
	}
}

func (a *CreateCampaignAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createCampaignActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	clientID := ctx.Auth.Extra["client-id"]
	if clientID == "" {
		return nil, fmt.Errorf("client ID is required")
	}

	if input.ListID == "" && len(input.SegmentIDs) == 0 {
		return nil, fmt.Errorf("either a list ID or segment IDs must be provided")
	}

	payload := map[string]interface{}{
		"Name":      input.Name,
		"Subject":   input.Subject,
		"FromName":  input.FromName,
		"FromEmail": input.FromEmail,
		"ReplyTo":   input.ReplyTo,
		"HtmlUrl":   input.HtmlURL,
	}

	if input.TextURL != "" {
		payload["TextUrl"] = input.TextURL
	}

	if input.ListID != "" {
		payload["ListIDs"] = []string{input.ListID}
	}

	if len(input.SegmentIDs) > 0 {
		payload["SegmentIDs"] = input.SegmentIDs
	}

	if input.InlineCSS != nil {
		payload["InlineCss"] = *input.InlineCSS
	}

	endpoint := fmt.Sprintf("campaigns/%s.json", clientID)

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

func (a *CreateCampaignAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateCampaignAction) SampleData() sdkcore.JSON {
	return "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1"
}

func (a *CreateCampaignAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateCampaignAction() sdk.Action {
	return &CreateCampaignAction{}
}
