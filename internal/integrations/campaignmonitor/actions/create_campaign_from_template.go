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

// Metadata returns metadata about the action
func (a *CreateCampaignFromTemplateAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_campaign_from_template",
		DisplayName:   "Create Campaign From Template",
		Description:   "Create a new draft campaign based on a template.",
		Type:          core.ActionTypeAction,
		Documentation: createCampaignTemplateDocs,
		Icon:          "mdi:email-edit",
		SampleOutput:  "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
		Settings:      core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CreateCampaignFromTemplateAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_campaign_from_template", "Create Campaign From Template")

	form.TextField("name", "Campaign Name").
		Placeholder("Enter campaign name").
		Required(true).
		HelpText("The name of your campaign.")

	form.TextField("subject", "Email Subject").
		Placeholder("Enter email subject").
		Required(true).
		HelpText("The subject line for your campaign.")

	form.TextField("fromName", "From Name").
		Placeholder("Enter sender name").
		Required(true).
		HelpText("The name that will appear as the sender.")

	form.TextField("fromEmail", "From Email").
		Placeholder("Enter sender email").
		Required(true).
		HelpText("The email address that will appear as the sender.")

	form.TextField("replyTo", "Reply-To Email").
		Placeholder("Enter reply-to email").
		Required(true).
		HelpText("The email address that recipients can reply to.")

	// Note: This will have type errors, but we're ignoring shared errors as per the issue description
	// form.SelectField("listIds", "List").
	//	Placeholder("Select a list").
	//	Required(true).
	//	WithDynamicOptions(...).
	//	HelpText("The list to send the campaign to.")

	form.TextField("templateId", "Template ID").
		Placeholder("Enter template ID").
		Required(true).
		HelpText("The ID of the template to use.")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CreateCampaignFromTemplateAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CreateCampaignFromTemplateAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[createCampaignFromTemplateActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Ensure at least one of ListIDs or SegmentIDs is provided
	if len(input.ListIDs) == 0 && len(input.SegmentIDs) == 0 {
		return nil, fmt.Errorf("at least one list ID or segment ID must be provided")
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	clientID := authCtx.Extra["client-id"]
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
		authCtx.Extra["api-key"],
		clientID,
		endpoint,
		http.MethodPost,
		payload)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewCreateCampaignFromTemplateAction() sdk.Action {
	return &CreateCampaignFromTemplateAction{}
}
