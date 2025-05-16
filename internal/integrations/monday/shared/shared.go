package shared

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gookit/goutil/arrutil"
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"

	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

// #nosec
var tokenURL = "https://auth.monday.com/oauth2/token"

const baseURL = "https://api.monday.com/v2"

var (
	form = smartform.NewAuthForm("monday-auth", "Monday.com Oauth", smartform.AuthStrategyOAuth2)
	_    = form.
		OAuthField("oauth", "Monday.com Oauth").
		AuthorizationURL("https://auth.monday.com/oauth2/authorize").
		TokenURL("https://auth.monday.com/oauth2/token").
		Scopes([]string{}).
		Build()
)

var SharedAuth = form.Build()

func MondayClient(ctx sdkcontext.BaseContext, query string) (map[string]interface{}, error) {
	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

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
	req.Header.Add("Authorization", token)
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

func GetWorkspaceProp(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getWorkspaces := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
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

		tokenSource := ctx.Auth().Token
		if tokenSource == nil {
			return nil, errors.New("missing authentication token")
		}
		token := tokenSource.AccessToken

		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Api-Version", "2023-07")
		req.Header.Set("Authorization", token)

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

	return form.SelectField("workspace_id", "Workspace").
		Placeholder("Select a workspace").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getWorkspaces)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText("Select a workspace")
}

func GetBoardProp(id, title, description string, form *smartform.FormBuilder) *smartform.FieldBuilder {
	getBoardID := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
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

		tokenSource := ctx.Auth().Token
		if tokenSource == nil {
			return nil, errors.New("missing authentication token")
		}
		token := tokenSource.AccessToken

		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Api-Version", "2023-07")
		req.Header.Set("Authorization", token)

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

	return form.SelectField(id, title).
		Placeholder(description).
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getBoardID)).
				WithSearchSupport().
				WithPagination(10).
				End().
				RefreshOn("workspace_id").
				GetDynamicSource(),
		).
		HelpText(description)
}

func GetGroupProp(id, title, description string, required bool, form *smartform.FormBuilder) *smartform.FieldBuilder {
	getGroups := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
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

		tokenSource := ctx.Auth().Token
		if tokenSource == nil {
			return nil, errors.New("missing authentication token")
		}
		token := tokenSource.AccessToken

		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Api-Version", "2023-07")
		req.Header.Set("Authorization", token)

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
				"value": input.ID,
				"label": input.Title,
			}, true
		})

		return ctx.Respond(items, len(items))
	}

	return form.SelectField(id, title).
		Placeholder(description).
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getGroups)).
				WithSearchSupport().
				WithPagination(10).
				End().
				RefreshOn("board_id").
				GetDynamicSource(),
		).
		HelpText(description)
}

var ColumnType = []*smartform.Option{
	{},
	{Value: "auto_number", Label: "Auto Number"},
	{Value: "board_relation", Label: "Board Relation"},
	{Value: "button", Label: "Button"},
	{Value: "checkbox", Label: "Checkbox"},
	{Value: "color_picker", Label: "Color Picker"},
	{Value: "country", Label: "Country"},
	{Value: "creation_log", Label: "Creation Log"},
	{Value: "date", Label: "Date"},
	{Value: "dependency", Label: "Dependency"},
	{Value: "doc", Label: "Doc"},
	{Value: "dropdown", Label: "Dropdown"},
	{Value: "email", Label: "Email"},
	{Value: "file", Label: "File"},
	{Value: "formula", Label: "Formula"},
	{Value: "hour", Label: "Hour"},
	{Value: "item_assignees", Label: "Item Assignees"},
	{Value: "item_id", Label: "Item ID"},
	{Value: "last_updated", Label: "Last Updated"},
	{Value: "link", Label: "Link"},
	{Value: "location", Label: "Location"},
	{Value: "long_text", Label: "Long Text"},
	{Value: "mirror", Label: "Mirror"},
	{Value: "name", Label: "Name"},
	{Value: "numbers", Label: "Numbers"},
	{Value: "phone", Label: "Phone"},
	{Value: "people", Label: "People"},
	{Value: "progress", Label: "Progress"},
	{Value: "rating", Label: "Rating"},
	{Value: "status", Label: "Status"},
	{Value: "subtasks", Label: "Subtasks"},
	{Value: "tags", Label: "Tags"},
	{Value: "team", Label: "Team"},
	{Value: "text", Label: "Text"},
	{Value: "timeline", Label: "Timeline"},
	{Value: "time_tracking", Label: "Time Tracking"},
	{Value: "vote", Label: "Vote"},
	{Value: "week", Label: "Week"},
	{Value: "world_clock", Label: "World Clock"},
	{Value: "unsupported", Label: "Unsupported"},
}
