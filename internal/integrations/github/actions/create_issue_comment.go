package actions

import (
	"errors"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/github/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createIssueCommentActionProps struct {
	Repository  string `json:"repository"`
	Body        string `json:"body"`
	IssueNumber string `json:"issue_number"`
}

type CreateIssueCommentAction struct{}

// Metadata returns metadata about the action
func (a *CreateIssueCommentAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_issue_comment",
		DisplayName:   "Create Issue Comment",
		Description:   "Create Issue Comment: Adds a comment to an existing GitHub issue, allowing for status updates, additional information, or discussion without creating a new issue.",
		Type:          core.ActionTypeAction,
		Documentation: createIssueCommentDocs,
		SampleOutput: map[string]any{
			"node": map[string]any{
				"id":        "MDExOlB1bGxSZXF1ZXN0Q29tbWVudDUxMjA0Mzc3NA==",
				"body":      "This is a comment",
				"createdAt": "2022-01-01T00:00:00Z",
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CreateIssueCommentAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_issue_comment", "Create Issue Comment")
	shared.RegisterRepositoryProps(form)
	shared.RegisterIssuesProps(form)
	form.TextareaField("body", "Comment").
		Placeholder("Enter your comment").
		Required(true).
		HelpText("The text content of your comment")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CreateIssueCommentAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CreateIssueCommentAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[createIssueCommentActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the token source from the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	mutation := fmt.Sprintf(`
			mutation {
			  addComment(input: { subjectId: "%s", body: "%s" }) {
				commentEdge {
				  node {
					id
					body
					createdAt
				  }
				}
			  }
			}`, input.IssueNumber, input.Body)

	response, err := shared.GithubGQL(authCtx.Token.AccessToken, mutation)
	if err != nil {
		return nil, errors.New("error making graphQL request")
	}

	comment, ok := response["data"].(map[string]interface{})["addComment"].(map[string]interface{})["commentEdge"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to extract comment from response")
	}

	return comment, nil
}

func NewCreateIssueCommentAction() sdk.Action {
	return &CreateIssueCommentAction{}
}
