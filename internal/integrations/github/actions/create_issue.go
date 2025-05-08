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

type createIssueActionProps struct {
	Repository  string `json:"repository"`
	Body        string `json:"body"`
	IssueNumber string `json:"issue_number"`
}

type CreateIssueAction struct{}

// Metadata returns metadata about the action
func (a *CreateIssueAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_issue",
		DisplayName:   "Create Issue",
		Description:   "Create Issue: Automatically generates a new issue in your project management tool (e.g., Jira, Trello) based on specific conditions or triggers, ensuring timely and organized tracking of tasks and projects.",
		Type:          core.ActionTypeAction,
		Documentation: createIssueDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CreateIssueAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_issue", "Create Issue")

	// Define the getRepositories function
	getRepositories := func(ctx sdkcontext.DynamicFieldContext) (*core.DynamicOptionsResponse, error) {
		return shared.GetRepositories(ctx)
	}

	// Define the getIssues function
	getIssues := func(ctx sdkcontext.DynamicFieldContext) (*core.DynamicOptionsResponse, error) {
		return shared.GetIssues(ctx)
	}

	// Add repository field
	form.SelectField("repository", "Repository").
		Placeholder("Select a repository").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getRepositories)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText("The repository to create the issue in.")

	// Add issue number field
	form.SelectField("issue_number", "Issue Number").
		Placeholder("Select an issue").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getIssues)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText("The issue to comment on.")

	// Add body field
	form.TextareaField("body", "Comment").
		Placeholder("Enter your comment").
		Required(true).
		HelpText("Issue comment")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CreateIssueAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CreateIssueAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[createIssueActionProps](ctx)
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

	response, err := shared.GithubGQL(authCtx.AccessToken, mutation)
	if err != nil {
		return nil, errors.New("error making graphQL request")
	}

	issue, ok := response["data"].(map[string]interface{})["addComment"].(map[string]interface{})["commentEdge"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to extract issue from response")
	}

	return issue, nil
}

func NewCreateIssueAction() sdk.Action {
	return &CreateIssueAction{}
}
