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

type createSpaceProps struct {
	WorkspaceID string `json:"workspace-id"`
	Name        string `json:"name"`
	Private     bool   `json:"private"`
}

type CreateSpaceOperation struct{}

// Metadata returns metadata about the action
func (o *CreateSpaceOperation) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_space",
		DisplayName:   "Create Space",
		Description:   "Creates a new space in a specified ClickUp workspace.",
		Type:          core.ActionTypeAction,
		Documentation: createSpaceDocs,
		Icon:          "material-symbols:create-new-folder",
		SampleOutput: map[string]any{
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
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (o *CreateSpaceOperation) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_space", "Create Space")

	shared.RegisterWorkSpaceInput(form, "Workspaces", "select a workspace", true)

	form.TextField("name", "name").
		Placeholder("Space Name").
		HelpText("The name of the space").
		Required(true)

	form.CheckboxField("private", "private").
		Placeholder("Private").
		HelpText("Whether the space is private").
		Required(false)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (o *CreateSpaceOperation) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (o *CreateSpaceOperation) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createSpaceProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	accessToken := authCtx.Token.AccessToken

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

func NewCreateSpaceOperation() sdk.Action {
	return &CreateSpaceOperation{}
}
