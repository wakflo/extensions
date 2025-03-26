package actions

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/jiracloudsoftware/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type transitionIssueActionProps struct {
	ProjectID    string `json:"projectId"`
	IssueID      string `json:"issueId"`
	TransitionID string `json:"transitionId"`
	Comment      string `json:"comment,omitempty"`
}

type TransitionIssueAction struct{}

func (a *TransitionIssueAction) Name() string {
	return "Transition Issue"
}

func (a *TransitionIssueAction) Description() string {
	return "Move a Jira issue to a different status"
}

func (a *TransitionIssueAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *TransitionIssueAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &transitionIssueDocs,
	}
}

func (a *TransitionIssueAction) Icon() *string {
	icon := "mdi:arrow-decision"
	return &icon
}

func (a *TransitionIssueAction) Properties() map[string]*sdkcore.AutoFormSchema {
	getTransitions := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		auth := ctx.Auth.Extra["email"] + ":" + ctx.Auth.Extra["api-token"]
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
		authHeader := "Basic " + encodedAuth

		input := sdk.DynamicInputToType[struct {
			ProjectID string `json:"projectId"`
			IssueID   string `json:"issueId"`
		}](ctx)

		if input.IssueID == "" {
			return ctx.Respond([]map[string]any{}, 0)
		}

		baseAPI := ctx.Auth.Extra["instance-url"] + "/rest/api/3/issue/" + input.IssueID + "/transitions"

		req, err := http.NewRequest(http.MethodGet, baseAPI, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", authHeader)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		client := &http.Client{}
		rsp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer rsp.Body.Close()

		responseBytes, err := io.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}

		var transitionResponse struct {
			Transitions []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				To   struct {
					Name string `json:"name"`
				} `json:"to"`
			} `json:"transitions"`
		}

		err = json.Unmarshal(responseBytes, &transitionResponse)
		if err != nil {
			return nil, err
		}

		items := make([]map[string]any, 0)
		for _, transition := range transitionResponse.Transitions {
			items = append(items, map[string]any{
				"id":   transition.ID,
				"name": transition.Name + " → " + transition.To.Name,
			})
		}

		return ctx.Respond(items, len(items))
	}

	return map[string]*sdkcore.AutoFormSchema{
		"projectId": shared.GetProjectsInput(),
		"issueId":   shared.GetIssuesInput(),
		"transitionId": autoform.NewDynamicField(sdkcore.String).
			SetDisplayName("Transition").
			SetDescription("Select a status transition").
			SetDynamicOptions(&getTransitions).
			SetRequired(true).Build(),
		"comment": autoform.NewLongTextField().
			SetDisplayName("Comment").
			SetDescription("Add a comment explaining this transition (optional)").
			SetRequired(false).Build(),
	}
}

func (a *TransitionIssueAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[transitionIssueActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	email := ctx.Auth.Extra["email"]
	apiToken := ctx.Auth.Extra["api-token"]
	instanceURL := ctx.Auth.Extra["instance-url"]

	// Build the request payload
	requestBody := map[string]interface{}{
		"transition": map[string]string{
			"id": input.TransitionID,
		},
	}

	// Add comment if provided
	if input.Comment != "" {
		requestBody["update"] = map[string]interface{}{
			"comment": []map[string]interface{}{
				{
					"add": map[string]interface{}{
						"body": map[string]interface{}{
							"type":    "doc",
							"version": 1,
							"content": []map[string]interface{}{
								{
									"type": "paragraph",
									"content": []map[string]interface{}{
										{
											"type": "text",
											"text": input.Comment,
										},
									},
								},
							},
						},
					},
				},
			},
		}
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	response, err := shared.JiraRequest(
		email,
		apiToken,
		instanceURL+"/rest/api/3/issue/"+input.IssueID+"/transitions",
		"POST",
		"Issue transitioned successfully",
		jsonBody,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *TransitionIssueAction) Auth() *sdk.Auth {
	return nil
}

func (a *TransitionIssueAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"Result": "Issue transitioned successfully",
	}
}

func (a *TransitionIssueAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewTransitionIssueAction() sdk.Action {
	return &TransitionIssueAction{}
}
