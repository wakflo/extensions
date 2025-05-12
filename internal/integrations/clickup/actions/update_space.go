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

type updateSpaceProps struct {
	SpaceID           string `json:"space-id"`
	SpaceName         string `json:"space-name"`
	MultipleAssignees bool   `json:"multiple-assignees"`
	Tags              bool   `json:"tags"`
	CustomFields      bool   `json:"custom-fields"`
}

type UpdateSpaceOperation struct{}

func (o *UpdateSpaceOperation) Name() string {
	return "Update Space"
}

func (o *UpdateSpaceOperation) Description() string {
	return "Updates an existing space in ClickUp with modified details."
}

func (o *UpdateSpaceOperation) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (o *UpdateSpaceOperation) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &updateSpaceDocs,
	}
}

func (o *UpdateSpaceOperation) Icon() *string {
	icon := "material-symbols:space-dashboard-edit"
	return &icon
}

func (o *UpdateSpaceOperation) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"workspace-id": shared.RegisterWorkSpaceInput("Workspaces", "select a workspace", true),
		"space-id":     shared.RegisterSpacesInput("Spaces", "select a space", true),
		"space-name": autoform.NewShortTextField().
			SetDisplayName("Update space Name").
			SetDescription("The space name to update").
			SetRequired(false).
			Build(),
		"multiple-assignees": autoform.NewBooleanField().
			SetDisplayName("Multiple Assignees").
			SetDescription("Enable multiple assignees").
			SetRequired(false).
			SetDefaultValue(true).
			Build(),
		"custom-fields": autoform.NewBooleanField().
			SetDisplayName("Custom fields").
			SetDescription("Enable custom fields").
			SetRequired(false).
			SetDefaultValue(true).
			Build(),
		"tags": autoform.NewBooleanField().
			SetDisplayName("Tags").
			SetDescription("Enable tags").
			SetRequired(false).
			SetDefaultValue(true).
			Build(),
	}
}

func (o *UpdateSpaceOperation) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	accessToken := ctx.Auth.AccessToken

	input, err := sdk.InputToTypeSafely[updateSpaceProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}
	reqURL := shared.BaseURL + "/v2/space/" + input.SpaceID
	data := []byte(fmt.Sprintf(`{
		"name": "%s",
		"multiple_assignees": %t,
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
				"enabled": %t
			},
			"time_estimates": {
				"enabled": true
			},
			"checklists": {
				"enabled": true
			},
			"custom_fields": {
				"enabled": %t
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
	}`, input.SpaceName, input.MultipleAssignees, input.Tags, input.CustomFields))
	req, err := http.NewRequest(http.MethodPut, reqURL, bytes.NewBuffer(data))
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
		"Result": "Space updated Successfully",
	}, nil
}

func (o *UpdateSpaceOperation) Auth() *sdk.Auth {
	return nil
}

func (o *UpdateSpaceOperation) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":      "123456",
		"name":    "Updated Space",
		"private": true,
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

func (o *UpdateSpaceOperation) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewUpdateSpaceOperation() sdk.Action {
	return &UpdateSpaceOperation{}
}
