package actions

import (
	"errors"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/github/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getIssueActionProps struct {
	Repository  string `json:"repository"`
	IssueNumber int    `json:"issue_number"`
}

type GetIssueAction struct{}

func (a *GetIssueAction) Name() string {
	return "Get Issue"
}

func (a *GetIssueAction) Description() string {
	return "Retrieves an issue from a specified issue tracking system or platform, allowing you to incorporate issue data into your automated workflows."
}

func (a *GetIssueAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetIssueAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getIssueDocs,
	}
}

func (a *GetIssueAction) Icon() *string {
	return nil
}

func (a *GetIssueAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"repository": shared.GetRepositoryInput(),
		"issue_number": autoform.NewNumberField().
			SetDisplayName("Issue Number").
			SetDescription("The issue number").
			SetRequired(true).
			Build(),
	}
}

func (a *GetIssueAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getIssueActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		query {
            node(id: "%s") {
				... on Repository {
				    issue(number:%d ){
					    title
						body
						createdAt
						updatedAt
						number
					author {
					  login
					}
					assignees(first:10){
					    nodes {
						     name
							 login
					   }
					}
				}
			}
		}
}`, input.Repository, input.IssueNumber)

	response, err := shared.GithubGQL(ctx.Auth.AccessToken, query)
	if err != nil {
		return nil, errors.New("error making graphQL request")
	}

	issue, ok := response["data"].(map[string]interface{})["node"].(map[string]interface{})["issue"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to extract issue from response")
	}

	return issue, nil
}

func (a *GetIssueAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetIssueAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *GetIssueAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetIssueAction() sdk.Action {
	return &GetIssueAction{}
}
