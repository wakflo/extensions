package actions

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type updateSpaceProps struct {
	SpaceID           string `json:"space-id"`
	SpaceName         string `json:"space-name"`
	MultipleAssignees bool   `json:"multiple-assignees"`
	Tags              bool   `json:"tags"`
	CustomFields      bool   `json:"custom-fields"`
}

type UpdateSpaceOperation struct{}

// Metadata returns metadata about the action
func (o *UpdateSpaceOperation) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "update_space",
		DisplayName:   "Update Space",
		Description:   "Updates an existing space in ClickUp with modified details.",
		Type:          core.ActionTypeAction,
		Documentation: updateSpaceDocs,
		Icon:          "material-symbols:space-dashboard-edit",
		SampleOutput: map[string]any{
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
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (o *UpdateSpaceOperation) Properties() *smartform.FormSchema {
	form := smartform.NewForm("update_space", "Update Space")

	shared.RegisterWorkSpaceInput(form, "Workspaces", "select a workspace", true)

	shared.RegisterSpacesInput(form, "Spaces", "select a space", true)

	form.TextField("space-name", "space-name").
		Placeholder("Update space Name").
		HelpText("The space name to update").
		Required(false)

	form.CheckboxField("multiple-assignees", "multiple-assignees").
		Placeholder("Multiple Assignees").
		HelpText("Enable multiple assignees").
		DefaultValue(true).
		Required(false)

	form.CheckboxField("custom-fields", "custom-fields").
		Placeholder("Custom fields").
		HelpText("Enable custom fields").
		DefaultValue(true).
		Required(false)

	form.CheckboxField("tags", "tags").
		Placeholder("Tags").
		HelpText("Enable tags").
		DefaultValue(true).
		Required(false)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (o *UpdateSpaceOperation) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (o *UpdateSpaceOperation) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateSpaceProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	accessToken := authCtx.Token.AccessToken
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

func NewUpdateSpaceOperation() sdk.Action {
	return &UpdateSpaceOperation{}
}
