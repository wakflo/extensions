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

type createIssueActionProps struct {
	Repository string `json:"repository"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	Labels     string `json:"labels"`
}

type CreateIssueAction struct{}

// Metadata returns metadata about the action
func (a *CreateIssueAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_issue",
		DisplayName:   "Create Issue",
		Description:   "Create Issue: Creates a new issue in a GitHub repository with customizable title, body, and labels.",
		Type:          core.ActionTypeAction,
		Documentation: createIssueDocs,
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

func (a *CreateIssueAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_issue", "Create Issue")

	shared.RegisterRepositoryProps(form)

	shared.RegisterLabelProps(form)

	form.TextField("title", "Issue Title").
		Placeholder("Enter an issue title").
		Required(true).
		HelpText("The title of the issue")

	form.TextareaField("body", "Description").
		Placeholder("Enter a description").
		Required(false).
		HelpText("Issue description")

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

	response, err := shared.GithubGQL(authCtx.Token.AccessToken, mutation)
	if err != nil {
		return nil, fmt.Errorf("error making graphQL request: %w", err)
	}

	return response, nil
}

func NewCreateIssueAction() sdk.Action {
	return &CreateIssueAction{}
}
