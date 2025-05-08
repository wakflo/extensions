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

type unlockIssueActionProps struct {
	Repository  string `json:"repository"`
	IssueNumber string `json:"issue_number"`
}

type UnlockIssueAction struct{}

// Metadata returns metadata about the action
func (a *UnlockIssueAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "unlock_issue",
		DisplayName:   "Unlock Issue",
		Description:   "Unlock Issue: Manually unlock an issue in your project management tool, allowing team members to view and work on it again.",
		Type:          core.ActionTypeAction,
		Documentation: unlockIssueDocs,
		SampleOutput: map[string]any{
			"locked": false,
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *UnlockIssueAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("unlock_issue", "Unlock Issue")

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
		HelpText("The repository to unlock the issue in.")

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
		HelpText("The issue to unlock.")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *UnlockIssueAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *UnlockIssueAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[unlockIssueActionProps](ctx)
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
		unlockLockable(input: { lockableId: "%s" }) {
			    unlockedRecord{
    				locked
  				}
	    }
	}`, input.IssueNumber)

	response, err := shared.GithubGQL(authCtx.AccessToken, mutation)
	if err != nil {
		return nil, errors.New("error making graphQL request")
	}

	issue, ok := response["data"].(map[string]interface{})["unlockLockable"].(map[string]interface{})["unlockedRecord"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to extract issue from response")
	}

	return issue, nil
}

func NewUnlockIssueAction() sdk.Action {
	return &UnlockIssueAction{}
}
