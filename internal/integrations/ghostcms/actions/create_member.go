package actions

import (
	_ "embed"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/ghostcms/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createMemberActionProps struct {
	Email       string   `json:"email"`
	Name        string   `json:"name,omitempty"`
	Note        string   `json:"note,omitempty"`
	Labels      []string `json:"labels,omitempty"`
	Newsletters []string `json:"newsletters,omitempty"`
	Subscribed  bool     `json:"subscribed"`
	Comped      bool     `json:"comped"`
}

type CreateMemberAction struct{}

func (a *CreateMemberAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_member",
		DisplayName:   "Create Member",
		Description:   "Creates a new member in your Ghost publication",
		Type:          core.ActionTypeAction,
		Documentation: createMemberDocs,
		SampleOutput: map[string]any{
			"id":         "63f8d9a0e4b0d7001c3e3b1b",
			"uuid":       "b2c3d4e5-f6a7-8901-bcde-f23456789012",
			"email":      "john.doe@example.com",
			"name":       "John Doe",
			"note":       "VIP customer",
			"status":     "free",
			"subscribed": true,
			"created_at": "2024-02-24T10:30:00.000Z",
			"updated_at": "2024-02-24T10:30:00.000Z",
			"labels": []map[string]any{
				{
					"id":   "1",
					"name": "Premium",
					"slug": "premium",
				},
			},
			"email_count":        0,
			"email_opened_count": 0,
			"email_open_rate":    0,
		},
		Settings: core.ActionSettings{},
	}
}

func (a *CreateMemberAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_member", "Create Member")

	form.EmailField("email", "Email").
		Required(true).
		HelpText("The email address of the member")

	form.TextField("name", "Name").
		Required(false).
		HelpText("The full name of the member")

	form.TextareaField("note", "Note").
		Required(false).
		HelpText("An internal note about the member")

	form.ArrayField("labels", "Labels").
		Required(false).
		HelpText("Labels to assign to the member")

	form.ArrayField("newsletters", "Newsletters").
		Required(false).
		HelpText("Newsletter IDs to subscribe the member to")

	form.CheckboxField("subscribed", "Subscribed").
		Required(false).
		HelpText("Whether the member is subscribed to newsletters").
		DefaultValue(true)

	form.CheckboxField("comped", "Complimentary Subscription").
		Required(false).
		HelpText("Whether to give the member a complimentary paid subscription").
		DefaultValue(false)

	schema := form.Build()
	return schema
}

func (a *CreateMemberAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *CreateMemberAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createMemberActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	client, err := shared.NewGhostClient(authCtx.Extra)
	if err != nil {
		return nil, err
	}

	// Build the member object
	member := map[string]interface{}{
		"email":      input.Email,
		"subscribed": input.Subscribed,
		"comped":     input.Comped,
	}

	// Add optional fields
	if input.Name != "" {
		member["name"] = input.Name
	}
	if input.Note != "" {
		member["note"] = input.Note
	}

	// Handle labels
	if len(input.Labels) > 0 {
		labels := make([]map[string]string, len(input.Labels))
		for i, label := range input.Labels {
			labels[i] = map[string]string{"name": label}
		}
		member["labels"] = labels
	}

	// Handle newsletters
	if len(input.Newsletters) > 0 {
		newsletters := make([]map[string]string, len(input.Newsletters))
		for i, newsletter := range input.Newsletters {
			newsletters[i] = map[string]string{"id": newsletter}
		}
		member["newsletters"] = newsletters
	}

	// Create the member
	response, err := client.Post("/members/", map[string]interface{}{
		"members": []interface{}{member},
	})
	if err != nil {
		return nil, err
	}

	// Extract the created member
	if members, ok := response["members"].([]interface{}); ok && len(members) > 0 {
		return members[0], nil
	}

	return response, nil
}

func NewCreateMemberAction() sdk.Action {
	return &CreateMemberAction{}
}
