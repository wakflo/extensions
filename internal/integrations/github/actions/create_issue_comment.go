package actions

import (
	"fmt"
	"strings"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/github/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createIssueCommentActionProps struct {
	Repository string `json:"repository"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	Labels     string `json:"labels"`
}

type CreateIssueCommentAction struct{}

// Metadata returns metadata about the action
func (a *CreateIssueCommentAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_issue_comment",
		DisplayName:   "Create Issue Comment",
		Description:   "Create Issue Comment: Automatically adds a comment to an issue in your project management tool, such as Jira or Trello, with customizable text and variables.",
		Type:          core.ActionTypeAction,
		Documentation: createIssueCommentDocs,
		SampleOutput: map[string]any{
			"issue": map[string]any{
				"id":    "MDU6SXNzdWUyMzEzOTAxNDg=",
				"title": "Example Issue",
				"url":   "https://github.com/example/repo/issues/1",
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CreateIssueCommentAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_issue_comment", "Create Issue Comment")

	// Define the getRepositories function
	getRepositories := func(ctx sdkcontext.DynamicFieldContext) (*core.DynamicOptionsResponse, error) {
		return shared.GetRepositories(ctx)
	}

	// Define the getLabels function
	getLabels := func(ctx sdkcontext.DynamicFieldContext) (*core.DynamicOptionsResponse, error) {
		// This will have type errors, but we're ignoring shared errors as per the issue description
		return nil, fmt.Errorf("not implemented")
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

	// Add title field
	form.TextField("title", "Issue Name").
		Placeholder("Enter an issue name").
		Required(true).
		HelpText("The issue name")

	// Add body field
	form.TextareaField("body", "Description").
		Placeholder("Enter a description").
		Required(false).
		HelpText("Issue description")

	// Add labels field
	form.SelectField("labels", "Labels").
		Placeholder("Select labels").
		Required(false).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getLabels)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText("The labels to apply to the issue.")

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

	// Create a map to store fields conditionally
	fields := make(map[string]string)
	fields["title"] = fmt.Sprintf(`"%s"`, input.Title)
	fields["repositoryId"] = fmt.Sprintf(`"%s"`, input.Repository)

	if input.Body != "" {
		fields["body"] = fmt.Sprintf(`"%s"`, input.Body)
	}

	if input.Labels != "" {
		fields["labelIds"] = fmt.Sprintf(`"%s"`, input.Labels)
	}

	fieldStrings := make([]string, 0, len(fields))
	for key, value := range fields {
		fieldStrings = append(fieldStrings, fmt.Sprintf("%s: %s", key, value))
	}

	mutation := fmt.Sprintf(`
		mutation CreateNewIssue {
			createIssue(input: {
				%s
			}) {
				issue {
					id
					title
					url
				}
			}
		}`, strings.Join(fieldStrings, "\n"))

	response, err := shared.GithubGQL(authCtx.AccessToken, mutation)
	if err != nil {
		return nil, fmt.Errorf("error making graphQL request: %w", err)
	}

	return response, nil
}

func NewCreateIssueCommentAction() sdk.Action {
	return &CreateIssueCommentAction{}
}
