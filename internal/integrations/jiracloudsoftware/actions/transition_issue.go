package actions

import (
	"encoding/json"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/jiracloudsoftware/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type transitionIssueActionProps struct {
	ProjectID    string `json:"projectId"`
	IssueID      string `json:"issueId"`
	TransitionID string `json:"transitionId"`
	Comment      string `json:"comment,omitempty"`
}

type TransitionIssueAction struct{}

// Metadata returns metadata about the action
func (a *TransitionIssueAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "transition_issue",
		DisplayName:   "Transition Issue",
		Description:   "Move a Jira issue to a different status",
		Type:          core.ActionTypeAction,
		Documentation: transitionIssueDocs,
		Icon:          "mdi:arrow-decision",
		SampleOutput: map[string]any{
			"Result": "Issue transitioned successfully",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *TransitionIssueAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("transition_issue", "Transition Issue")

	// Register project selection field
	shared.RegisterProjectsProps(form)

	// Register issue selection field
	shared.RegisterIssuesProps(form)

	shared.RegisterTransitionsProps(form)

	form.TextField("comment", "Comment").
		Required(false).
		HelpText("Add a comment explaining this transition (optional)")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *TransitionIssueAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *TransitionIssueAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[transitionIssueActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	email := authCtx.Extra["email"]
	apiToken := authCtx.Extra["api-token"]
	instanceURL := authCtx.Extra["instance-url"]

	// Build the request payload
	requestBody := map[string]interface{}{
		"transition": map[string]string{
			"id": input.TransitionID,
		},
	}

	// Add comment if provided
	if input.Comment != "" {
		requestBody["update"] = map[string]interface{}{
			"comment": []map[string]interface{}{
				{
					"add": map[string]interface{}{
						"body": map[string]interface{}{
							"type":    "doc",
							"version": 1,
							"content": []map[string]interface{}{
								{
									"type": "paragraph",
									"content": []map[string]interface{}{
										{
											"type": "text",
											"text": input.Comment,
										},
									},
								},
							},
						},
					},
				},
			},
		}
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	response, err := shared.JiraRequest(
		email,
		apiToken,
		instanceURL+"/rest/api/3/issue/"+input.IssueID+"/transitions",
		"POST",
		"Issue transitioned successfully",
		jsonBody,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewTransitionIssueAction() sdk.Action {
	return &TransitionIssueAction{}
}
