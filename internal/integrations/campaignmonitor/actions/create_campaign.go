package actions

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/campaignmonitor/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
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

// Metadata returns metadata about the action
func (a *CreateCampaignAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_campaign",
		DisplayName:   "Create Campaign",
		Description:   "Create a new draft campaign ready to be tested or sent.",
		Type:          core.ActionTypeAction,
		Documentation: createCampaignDocs,
		Icon:          "mdi:email-newsletter",
		SampleOutput:  "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
		Settings:      core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CreateCampaignAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_campaign", "Create Campaign")

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

	form.TextField("htmlUrl", "HTML URL").
		Placeholder("Enter HTML URL").
		Required(true).
		HelpText("The URL of the HTML content for your campaign.")

	form.TextField("textUrl", "Text URL").
		Placeholder("Enter text URL").
		Required(false).
		HelpText("The URL of the text content for your campaign (optional).")

	// Note: This will have type errors, but we're ignoring shared errors as per the issue description
	// form.SelectField("listId", "List").
	//	Placeholder("Select a list").
	//	Required(true).
	//	WithDynamicOptions(...).
	//	HelpText("The list to send the campaign to.")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CreateCampaignAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CreateCampaignAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[createCampaignActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	clientID := authCtx.Extra["client-id"]
	if clientID == "" {
		return nil, errors.New("client ID is required")
	}

	if input.ListID == "" && len(input.SegmentIDs) == 0 {
		return nil, errors.New("either a list ID or segment IDs must be provided")
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

func NewCreateCampaignAction() sdk.Action {
	return &CreateCampaignAction{}
}
