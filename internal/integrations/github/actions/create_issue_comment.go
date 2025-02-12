package actions

import (
	"fmt"
	"strings"

	"github.com/wakflo/extensions/internal/integrations/github/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createIssueCommentActionProps struct {
	Repository string `json:"repository"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	Labels     string `json:"labels"`
}

type CreateIssueCommentAction struct{}

func (a *CreateIssueCommentAction) Name() string {
	return "Create Issue Comment"
}

func (a *CreateIssueCommentAction) Description() string {
	return "Create Issue Comment: Automatically adds a comment to an issue in your project management tool, such as Jira or Trello, with customizable text and variables."
}

func (a *CreateIssueCommentAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateIssueCommentAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createIssueCommentDocs,
	}
}

func (a *CreateIssueCommentAction) Icon() *string {
	return nil
}

func (a *CreateIssueCommentAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"repository": shared.GetRepositoryInput(),
		"title": autoform.NewShortTextField().
			SetDisplayName("Issue Name").
			SetDescription("The issue name").
			SetRequired(true).
			Build(),
		"body": autoform.NewLongTextField().
			SetDisplayName("Description").
			SetDescription("Issue description").
			Build(),
		"labels": shared.GetLabelInput(),
	}
}

func (a *CreateIssueCommentAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createIssueCommentActionProps](ctx.BaseContext)
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

	response, err := shared.GithubGQL(ctx.Auth.AccessToken, mutation)
	if err != nil {
		return nil, fmt.Errorf("error making graphQL request: %w", err)
	}

	return response, nil
}

func (a *CreateIssueCommentAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateIssueCommentAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateIssueCommentAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateIssueCommentAction() sdk.Action {
	return &CreateIssueCommentAction{}
}
