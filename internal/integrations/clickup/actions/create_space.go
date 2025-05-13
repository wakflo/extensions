package actions

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createSpaceProps struct {
	WorkspaceID string `json:"workspace-id"`
	Name        string `json:"name"`
	Private     bool   `json:"private"`
}

type CreateSpaceOperation struct{}

func (o *CreateSpaceOperation) Name() string {
	return "Create Space"
}

func (o *CreateSpaceOperation) Description() string {
	return "Creates a new space in a specified ClickUp workspace."
}

func (o *CreateSpaceOperation) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (o *CreateSpaceOperation) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createSpaceDocs,
	}
}

func (o *CreateSpaceOperation) Icon() *string {
	icon := "material-symbols:create-new-folder"
	return &icon
}

func (o *CreateSpaceOperation) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"workspace-id": shared.RegisterWorkSpaceInput("Workspaces", "select a workspace", true),
		"name": autoform.NewShortTextField().
			SetDisplayName("Space Name").
			SetDescription("The name of the space").
			SetRequired(true).
			Build(),
	}
}

func (o *CreateSpaceOperation) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[createSpaceProps](ctx.BaseContext)

	reqURL := shared.BaseURL + "/v2/team/" + input.WorkspaceID + "/space"
	data := []byte(fmt.Sprintf(`{
		"name": "%s",
		"multiple_assignees": true,
		"features": {
			"due_dates": {
				"enabled": true,
				"start_date": false,
				"remap_due_dates": true,
				"remap_closed_due_date": false
			},
			"time_tracking": {
				"enabled": false
			},
			"tags": {
				"enabled": true
			},
			"time_estimates": {
				"enabled": true
			},
			"checklists": {
				"enabled": true
			},
			"custom_fields": {
				"enabled": true
			},
			"remap_dependencies": {
				"enabled": true
			},
			"dependency_warning": {
				"enabled": true
			},
			"portfolios": {
				"enabled": true
			}
		}
	}`, input.Name))
	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", accessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(res)
	fmt.Println(string(body))

	return map[string]interface{}{
		"Result": "Space created Successfully",
	}, nil
}

func (o *CreateSpaceOperation) Auth() *sdk.Auth {
	return nil
}

func (o *CreateSpaceOperation) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":      "123456",
		"name":    "New Space",
		"private": false,
		"statuses": []map[string]any{
			{
				"id":     "st123",
				"status": "Open",
				"color":  "#d3d3d3",
			},
		},
		"multiple_assignees": true,
		"features": map[string]any{
			"due_dates":     true,
			"time_tracking": true,
			"tags":          true,
		},
	}
}

func (o *CreateSpaceOperation) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateSpaceOperation() sdk.Action {
	return &CreateSpaceOperation{}
}
