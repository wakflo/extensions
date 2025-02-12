package actions

import (
	"errors"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/github/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type unlockIssueActionProps struct {
	Repository  string `json:"repository"`
	IssueNumber string `json:"issue_number"`
}

type UnlockIssueAction struct{}

func (a *UnlockIssueAction) Name() string {
	return "Unlock Issue"
}

func (a *UnlockIssueAction) Description() string {
	return "Unlock Issue: Manually unlock an issue in your project management tool, allowing team members to view and work on it again."
}

func (a *UnlockIssueAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *UnlockIssueAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &unlockIssueDocs,
	}
}

func (a *UnlockIssueAction) Icon() *string {
	return nil
}

func (a *UnlockIssueAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"repository":   shared.GetRepositoryInput(),
		"issue_number": shared.GetIssuesInput(),
	}
}

func (a *UnlockIssueAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[unlockIssueActionProps](ctx.BaseContext)
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

	response, err := shared.GithubGQL(ctx.Auth.AccessToken, mutation)
	if err != nil {
		return nil, errors.New("error making graphQL request")
	}

	issue, ok := response["data"].(map[string]interface{})["unlockLockable"].(map[string]interface{})["unlockedRecord"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to extract issue from response")
	}

	return issue, nil
}

func (a *UnlockIssueAction) Auth() *sdk.Auth {
	return nil
}

func (a *UnlockIssueAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *UnlockIssueAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewUnlockIssueAction() sdk.Action {
	return &UnlockIssueAction{}
}
