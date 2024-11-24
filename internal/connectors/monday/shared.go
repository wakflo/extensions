package monday

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gookit/goutil/arrutil"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var (
	// #nosec
	tokenURL   = "https://auth.monday.com/oauth2/token"
	sharedAuth = autoform.NewOAuthField("https://auth.monday.com/oauth2/authorize", &tokenURL, []string{}).Build()
)

const baseURL = "https://api.monday.com/v2"

func mondayClient(accessToken, query string) (map[string]interface{}, error) {
	payload := map[string]string{
		"query": query,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Api-Version", "2023-07")
	req.Header.Add("Authorization", accessToken)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var mondayResponse map[string]interface{}

	err = json.Unmarshal(body, &mondayResponse)
	if err != nil {
		return nil, err
	}

	return mondayResponse, nil
}

func getWorkspaceInput() *sdkcore.AutoFormSchema {
	getWorkspaces := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		query := `{
		 workspaces {
		     id
		     name
		 }
		}`

		queryBody := map[string]string{
			"query": query,
		}
		jsonQuery, err := json.Marshal(queryBody)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(jsonQuery))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Api-Version", "2023-07")
		req.Header.Set("Authorization", ctx.Auth.AccessToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var response struct {
			Data struct {
				Workspaces []Workspace `json:"workspaces"`
			} `json:"data"`
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		teams := response.Data.Workspaces
		return ctx.Respond(teams, len(teams))
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Workspaces").
		SetDescription("Select a workspace").
		SetDynamicOptions(&getWorkspaces).
		SetRequired(true).Build()
}

func getBoardInput(title, description string) *sdkcore.AutoFormSchema {
	getBoardID := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		input := sdk.DynamicInputToType[struct {
			WorkspaceID string `json:"workspace_id"`
		}](ctx)

		query := fmt.Sprintf(`{
		     boards (workspace_ids: %s, order_by: created_at) {
						id
    					name
  				}
        }`, input.WorkspaceID)

		queryBody := map[string]string{
			"query": query,
		}
		jsonQuery, err := json.Marshal(queryBody)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(jsonQuery))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Api-Version", "2023-07")
		req.Header.Set("Authorization", ctx.Auth.AccessToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var response struct {
			Data struct {
				Boards []Board `json:"boards"`
			} `json:"data"`
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		boards := response.Data.Boards
		return ctx.Respond(boards, len(boards))
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(description).
		SetDynamicOptions(&getBoardID).
		SetRequired(true).Build()
}

func getGroupInput(title, description string, required bool) *sdkcore.AutoFormSchema {
	getGroups := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		input := sdk.DynamicInputToType[struct {
			BoardID string `json:"board_id"`
		}](ctx)

		query := fmt.Sprintf(`
    {
        boards(ids: %s){
            groups  {
                id
                title
            }
        }
    }`, input.BoardID)

		queryBody := map[string]string{
			"query": query,
		}
		jsonQuery, err := json.Marshal(queryBody)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(jsonQuery))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Api-Version", "2023-07")
		req.Header.Set("Authorization", ctx.Auth.AccessToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var response Response

		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}
		if len(response.Data.Boards) == 0 {
			return nil, errors.New("no boards found")
		}

		groups := response.Data.Boards[0].Groups

		items := arrutil.Map[Group, map[string]any](groups, func(input Group) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.ID,
				"name": input.Title,
			}, true
		})

		return ctx.Respond(items, len(items))
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(description).
		SetDynamicOptions(&getGroups).
		SetRequired(required).Build()
}

var columnType = []*sdkcore.AutoFormSchema{
	{Const: "auto_number", Title: "Auto Number"},
	{Const: "board_relation", Title: "Board Relation"},
	{Const: "button", Title: "Button"},
	{Const: "checkbox", Title: "Checkbox"},
	{Const: "color_picker", Title: "Color Picker"},
	{Const: "country", Title: "Country"},
	{Const: "creation_log", Title: "Creation Log"},
	{Const: "date", Title: "Date"},
	{Const: "dependency", Title: "Dependency"},
	{Const: "doc", Title: "Doc"},
	{Const: "dropdown", Title: "Dropdown"},
	{Const: "email", Title: "Email"},
	{Const: "file", Title: "File"},
	{Const: "formula", Title: "Formula"},
	{Const: "hour", Title: "Hour"},
	{Const: "item_assignees", Title: "Item Assignees"},
	{Const: "item_id", Title: "Item ID"},
	{Const: "last_updated", Title: "Last Updated"},
	{Const: "link", Title: "Link"},
	{Const: "location", Title: "Location"},
	{Const: "long_text", Title: "Long Text"},
	{Const: "mirror", Title: "Mirror"},
	{Const: "name", Title: "Name"},
	{Const: "numbers", Title: "Numbers"},
	{Const: "phone", Title: "Phone"},
	{Const: "people", Title: "People"},
	{Const: "progress", Title: "Progress"},
	{Const: "rating", Title: "Rating"},
	{Const: "status", Title: "Status"},
	{Const: "subtasks", Title: "Subtasks"},
	{Const: "tags", Title: "Tags"},
	{Const: "team", Title: "Team"},
	{Const: "text", Title: "Text"},
	{Const: "timeline", Title: "Timeline"},
	{Const: "time_tracking", Title: "Time Tracking"},
	{Const: "vote", Title: "Vote"},
	{Const: "week", Title: "Week"},
	{Const: "world_clock", Title: "World Clock"},
	{Const: "unsupported", Title: "Unsupported"},
}
