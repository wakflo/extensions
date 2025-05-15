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

type getIssueActionProps struct {
	Repository  string `json:"repository"`
	IssueNumber int    `json:"issue_number"`
}

type GetIssueAction struct{}

// Metadata returns metadata about the action
func (a *GetIssueAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_issue",
		DisplayName:   "Get Issue",
		Description:   "Retrieves an issue from a specified issue tracking system or platform, allowing you to incorporate issue data into your automated workflows.",
		Type:          core.ActionTypeAction,
		Documentation: getIssueDocs,
		SampleOutput: map[string]any{
			"title":     "Example Issue",
			"body":      "This is an example issue",
			"createdAt": "2023-01-01T00:00:00Z",
			"updatedAt": "2023-01-02T00:00:00Z",
			"number":    42,
			"author":    map[string]any{"login": "username"},
			"assignees": map[string]any{"nodes": []map[string]any{{"name": "User Name", "login": "username"}}},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GetIssueAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_issue", "Get Issue")

	shared.RegisterRepositoryProps(form)
	shared.RegisterIssuesProps(form)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetIssueAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GetIssueAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[getIssueActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the token source from the auth context
	authCtx, err := ctx.AuthContext()
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

	response, err := shared.GithubGQL(authCtx.AccessToken, query)
	if err != nil {
		return nil, errors.New("error making graphQL request")
	}

	issue, ok := response["data"].(map[string]interface{})["node"].(map[string]interface{})["issue"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to extract issue from response")
	}

	return issue, nil
}

func NewGetIssueAction() sdk.Action {
	return &GetIssueAction{}
}
