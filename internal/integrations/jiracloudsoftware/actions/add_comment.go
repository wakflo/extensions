package actions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/jiracloudsoftware/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type addCommentActionProps struct {
	ProjectID   string `json:"projectId"`
	IssueID     string `json:"issueId"`
	CommentText string `json:"commentText"`
}

type AddCommentAction struct{}

// Metadata returns metadata about the action
func (a *AddCommentAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "add_comment",
		DisplayName:   "Add Comment",
		Description:   "Add a comment to a Jira issue",
		Type:          core.ActionTypeAction,
		Documentation: addCommentDocs,
		Icon:          "mdi:comment-plus",
		SampleOutput: map[string]any{
			"id":   "12345",
			"self": "https://yourcompany.atlassian.net/rest/api/3/issue/PRJ-123/comment/12345",
			"body": map[string]any{
				"type":    "doc",
				"version": 1,
				"content": []map[string]any{
					{
						"type": "paragraph",
						"content": []map[string]any{
							{
								"type": "text",
								"text": "This is a sample comment.",
							},
						},
					},
				},
			},
			"author": map[string]any{
				"displayName":  "John Doe",
				"emailAddress": "john.doe@example.com",
			},
			"created": "2023-01-15T10:30:45.123+0000",
			"updated": "2023-01-15T10:30:45.123+0000",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *AddCommentAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("add_comment", "Add Comment")

	shared.RegisterProjectsProps(form)

	shared.RegisterIssuesProps(form)

	form.TextareaField("commentText", "Comment").
		Required(true).
		HelpText("Text of the comment to add")

	schema := form.Build()

	return schema
}

func (a *AddCommentAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *AddCommentAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[addCommentActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("/rest/api/3/issue/%s/comment", input.IssueID)

	instanceURL := authCtx.Extra["instance-url"] + endpoint

	payload := map[string]interface{}{
		"body": map[string]interface{}{
			"type":    "doc",
			"version": 1,
			"content": []map[string]interface{}{
				{
					"type": "paragraph",
					"content": []map[string]interface{}{
						{
							"text": input.CommentText,
							"type": "text",
						},
					},
				},
			},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := shared.JiraRequest(authCtx.Extra["email"], authCtx.Extra["api-token"], instanceURL, http.MethodPost, "Comment added successfully!", data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func NewAddCommentAction() sdk.Action {
	return &AddCommentAction{}
}
