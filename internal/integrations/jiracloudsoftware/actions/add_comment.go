package actions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/jiracloudsoftware/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type addCommentActionProps struct {
	ProjectID   string `json:"projectId"`
	IssueID     string `json:"issueId"`
	CommentText string `json:"commentText"`
}

type AddCommentAction struct{}

func (a *AddCommentAction) Name() string {
	return "Add Comment"
}

func (a *AddCommentAction) Description() string {
	return "Add a comment to a Jira issue"
}

func (a *AddCommentAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *AddCommentAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &addCommentDocs,
	}
}

func (a *AddCommentAction) Icon() *string {
	icon := "mdi:comment-plus"
	return &icon
}

func (a *AddCommentAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"projectId": shared.GetProjectsInput(),
		"issueId":   shared.GetIssuesInput(),
		"commentText": autoform.NewLongTextField().
			SetDisplayName("Comment").
			SetDescription("Text of the comment to add").
			SetRequired(true).Build(),
	}
}

func (a *AddCommentAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[addCommentActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("/rest/api/3/issue/%s/comment", input.IssueID)

	instanceURL := ctx.Auth.Extra["instance-url"] + endpoint

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

	resp, err := shared.JiraRequest(ctx.Auth.Extra["email"], ctx.Auth.Extra["api-token"], instanceURL, http.MethodPost, "Comment added successfully!", data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *AddCommentAction) Auth() *sdk.Auth {
	return nil
}

func (a *AddCommentAction) SampleData() sdkcore.JSON {
	return map[string]any{
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
	}
}

func (a *AddCommentAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewAddCommentAction() sdk.Action {
	return &AddCommentAction{}
}
