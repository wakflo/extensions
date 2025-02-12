package actions

import (
	"errors"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/github/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createIssueActionProps struct {
	Repository  string `json:"repository"`
	Body        string `json:"body"`
	IssueNumber string `json:"issue_number"`
}

type CreateIssueAction struct{}

func (a *CreateIssueAction) Name() string {
	return "Create Issue"
}

func (a *CreateIssueAction) Description() string {
	return "Create Issue: Automatically generates a new issue in your project management tool (e.g., Jira, Trello) based on specific conditions or triggers, ensuring timely and organized tracking of tasks and projects."
}

func (a *CreateIssueAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateIssueAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createIssueDocs,
	}
}

func (a *CreateIssueAction) Icon() *string {
	return nil
}

func (a *CreateIssueAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"repository":   shared.GetRepositoryInput(),
		"issue_number": shared.GetIssuesInput(),
		"body": autoform.NewLongTextField().
			SetDisplayName("Comment").
			SetDescription("Issue comment").
			Build(),
	}
}

func (a *CreateIssueAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createIssueActionProps](ctx.BaseContext)
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

	response, err := shared.GithubGQL(ctx.Auth.AccessToken, mutation)
	if err != nil {
		return nil, errors.New("error making graphQL request")
	}

	issue, ok := response["data"].(map[string]interface{})["addComment"].(map[string]interface{})["commentEdge"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to extract issue from response")
	}

	return issue, nil
}

func (a *CreateIssueAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateIssueAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateIssueAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateIssueAction() sdk.Action {
	return &CreateIssueAction{}
}
