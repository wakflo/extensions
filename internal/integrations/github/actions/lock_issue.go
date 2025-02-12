package actions

import (
	"errors"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/github/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type lockIssueActionProps struct {
	Repository  string `json:"repository"`
	LockReason  string `json:"lock_reason"`
	IssueNumber string `json:"issue_number"`
}

type LockIssueAction struct{}

func (a *LockIssueAction) Name() string {
	return "Lock Issue"
}

func (a *LockIssueAction) Description() string {
	return "Locks an issue in the workflow, preventing any further updates or changes until it is manually unlocked."
}

func (a *LockIssueAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *LockIssueAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &lockIssueDocs,
	}
}

func (a *LockIssueAction) Icon() *string {
	return nil
}

func (a *LockIssueAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"repository":   shared.GetRepositoryInput(),
		"issue_number": shared.GetIssuesInput(),
		"lock_reason": autoform.NewSelectField().
			SetDisplayName("Lock Reason").
			SetDescription("The reason for locking the issue").
			SetOptions(shared.LockIssueReason).
			SetRequired(true).
			Build(),
	}
}

func (a *LockIssueAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[lockIssueActionProps](ctx.BaseContext)
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

	response, err := shared.GithubGQL(ctx.Auth.AccessToken, mutation)
	if err != nil {
		return nil, errors.New("error making graphQL request")
	}

	issue, ok := response["data"].(map[string]interface{})["lockLockable"].(map[string]interface{})["lockedRecord"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to extract issue from response")
	}

	return issue, nil
}

func (a *LockIssueAction) Auth() *sdk.Auth {
	return nil
}

func (a *LockIssueAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *LockIssueAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewLockIssueAction() sdk.Action {
	return &LockIssueAction{}
}
