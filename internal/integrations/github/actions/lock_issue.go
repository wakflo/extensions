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

type lockIssueActionProps struct {
	Repository  string `json:"repository"`
	LockReason  string `json:"lock_reason"`
	IssueNumber string `json:"issue_number"`
}

type LockIssueAction struct{}

// Metadata returns metadata about the action
func (a *LockIssueAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "lock_issue",
		DisplayName:   "Lock Issue",
		Description:   "Locks an issue in the workflow, preventing any further updates or changes until it is manually unlocked.",
		Type:          core.ActionTypeAction,
		Documentation: lockIssueDocs,
		SampleOutput: map[string]any{
			"locked":           true,
			"activeLockReason": "OFF_TOPIC",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *LockIssueAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("lock_issue", "Lock Issue")

	shared.RegisterRepositoryProps(form)
	shared.RegisterIssuesProps(form)
	form.SelectField("lock_reason", "Lock Reason").
		Placeholder("Select a reason").
		Required(true).
		AddOptions([]*smartform.Option{
			{Value: "OFF_TOPIC", Label: "Off Topic"},
			{Value: "TOO_HEATED", Label: "Too Heated"},
			{Value: "RESOLVED", Label: "Resolved"},
			{Value: "SPAM", Label: "Spam"},
		}...).
		HelpText("The reason for locking the issue.")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *LockIssueAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *LockIssueAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[lockIssueActionProps](ctx)
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
			lockLockable(input: { lockableId: "%s", lockReason: %s }) {
			lockedRecord {
				locked
        		activeLockReason
			}
		}
	}`, input.IssueNumber, input.LockReason)

	response, err := shared.GithubGQL(authCtx.Token.AccessToken, mutation)
	if err != nil {
		return nil, errors.New("error making graphQL request")
	}

	issue, ok := response["data"].(map[string]interface{})["lockLockable"].(map[string]interface{})["lockedRecord"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to extract issue from response")
	}

	return issue, nil
}

func NewLockIssueAction() sdk.Action {
	return &LockIssueAction{}
}
