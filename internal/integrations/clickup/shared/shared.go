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
	fastshot "github.com/opus-domini/fast-shot"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"

	"github.com/wakflo/go-sdk/v2"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

var (
	form = smartform.NewAuthForm("clickup-auth", "ClickUp OAuth", smartform.AuthStrategyOAuth2)
	_    = form.
		OAuthField("oauth", "ClickUp OAuth").
		AuthorizationURL("https://app.clickup.com/api").
		TokenURL("https://api.clickup.com/api/v2/oauth/token").
		Scopes([]string{}).
		Build()
)

var ClickupSharedAuth = form.Build()

const BaseURL = "https://api.clickup.com/api"

func GetAllSpaces(accessToken, param string) (interface{}, error) {
	reqURL := BaseURL + "/v2/team/" + param + "/space"
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)

	query := req.URL.Query()
	query.Add("archived", "false")
	req.URL.RawQuery = query.Encode()

	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", accessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var respData SpacesResponse
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return nil, err
	}

	space := respData.Spaces

	return arrutil.Map[Space, map[string]any](space, func(input Space) (target map[string]any, find bool) {
		return map[string]any{
			"id":   input.ID,
			"name": input.Name,
		}, true
	}), nil
}

func GetData(accessToken, url string) (map[string]interface{}, error) {
	reqURL := BaseURL + url
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)

	query := req.URL.Query()
	req.URL.RawQuery = query.Encode()

	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", accessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var respData map[string]interface{}
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return nil, err
	}

	return respData, nil
}

func GetList(accessToken, listID string) (map[string]interface{}, error) {
	reqURL := BaseURL + "/v2/list/" + listID
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", accessToken)

	fmt.Println("Request URL:", req.URL.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
	fmt.Println(string(body))

	var respData map[string]interface{}
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return nil, err
	}
	return respData, nil
}

func SearchTask(accessToken, url string, searchQuery string) (map[string]interface{}, error) {
	fullURL := BaseURL + url

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("query", searchQuery)

	req.URL.RawQuery = query.Encode()

	req.Header.Add("Authorization", accessToken)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

func GetTeams(accessToken string) ([]Team, error) {
	url := BaseURL + "/v2/team"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get teams from ClickUp API")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsedResponse struct {
		Teams []Team `json:"teams"`
	}

	if errs := json.Unmarshal(body, &parsedResponse); errs != nil {
		return nil, err
	}

	return parsedResponse.Teams, nil
}

func RegisterWorkSpaceInput(form *smartform.FormBuilder, title string, desc string, required bool) {
	getWorkspaces := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		client := fastshot.NewClient(BaseURL).
			Auth().BearerToken(ctx.Auth().AccessToken).
			Header().
			AddAccept("application/json").
			Build()

		rsp, err := client.GET("/v2/team").Send()
		if err != nil {
			return nil, err
		}

		defer rsp.Body().Close()

		if rsp.Status().IsError() {
			return nil, errors.New(rsp.Status().Text())
		}

		defer rsp.Body().Close()
		byts, err := io.ReadAll(rsp.Body().Raw())
		if err != nil {
			return nil, err
		}

		var body TeamsResponse

		err = json.Unmarshal(byts, &body)
		if err != nil {
			return nil, err
		}
		return ctx.Respond(body.Teams, len(body.Teams))
	}

	form.SelectField("workspace-id", title).
		Placeholder("Select a workspace").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getWorkspaces)).
				// WithFieldReference("state", "state").
				WithSearchSupport().
				WithPagination(10).
				End().
				// RefreshOn("state").
				GetDynamicSource(),
		).
		HelpText(desc)
}

func RegisterSpacesInput(form *smartform.FormBuilder, title string, desc string, required bool) {
	getSpaces := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		input := sdk.DynamicInputToType[struct {
			WorkspaceID string `json:"workspace-id,omitempty"`
		}](ctx)

		client := fastshot.NewClient(BaseURL).
			Auth().BearerToken(ctx.Auth().AccessToken).
			Header().
			AddAccept("application/json").
			Build()

		rsp, err := client.GET(fmt.Sprintf("/v2/team/%s/space", input.WorkspaceID)).Send()
		if err != nil {
			return nil, err
		}

		defer rsp.Body().Close()

		if rsp.Status().IsError() {
			return nil, errors.New(rsp.Status().Text())
		}

		defer rsp.Body().Close()
		byts, err := io.ReadAll(rsp.Body().Raw())
		if err != nil {
			return nil, err
		}

		var body SpacesResponse

		err = json.Unmarshal(byts, &body)
		if err != nil {
			return nil, err
		}

		return ctx.Respond(body.Spaces, len(body.Spaces))
	}

	form.SelectField("space-id", title).
		Placeholder("Select a space").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getSpaces)).
				WithFieldReference("workspace-id", "workspace-id").
				WithSearchSupport().
				WithPagination(10).
				End().
				RefreshOn("workspace-id").
				GetDynamicSource(),
		).
		HelpText(desc)
}

func RegisterFoldersInput(form *smartform.FormBuilder, title string, desc string, required bool) {
	getFolders := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		input := sdk.DynamicInputToType[struct {
			SpaceID string `json:"space-id,omitempty"`
		}](ctx)

		client := fastshot.NewClient(BaseURL).
			Auth().BearerToken(ctx.Auth().AccessToken).
			Header().
			AddAccept("application/json").
			Build()

		rsp, err := client.GET(fmt.Sprintf("/v2/space/%s/folder", input.SpaceID)).Send()
		if err != nil {
			return nil, err
		}

		defer rsp.Body().Close()

		if rsp.Status().IsError() {
			return nil, errors.New(rsp.Status().Text())
		}

		defer rsp.Body().Close()
		byts, err := io.ReadAll(rsp.Body().Raw())
		if err != nil {
			return nil, err
		}

		var body FoldersResponse

		err = json.Unmarshal(byts, &body)
		if err != nil {
			return nil, err
		}

		return ctx.Respond(body.Folders, len(body.Folders))
	}

	form.SelectField("folder-id", title).
		Placeholder("Select a folder").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getFolders)).
				WithFieldReference("space-id", "space-id").
				WithSearchSupport().
				WithPagination(10).
				End().
				RefreshOn("space-id").
				GetDynamicSource(),
		).
		HelpText(desc)
}

func RegisterListsInput(form *smartform.FormBuilder, title string, desc string, required bool) {
	getLists := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		input := sdk.DynamicInputToType[struct {
			FolderID string `json:"folder-id,omitempty"`
		}](ctx)

		client := fastshot.NewClient(BaseURL).
			Auth().BearerToken(ctx.Auth().AccessToken).
			Header().
			AddAccept("application/json").
			Build()

		rsp, err := client.GET(fmt.Sprintf("/v2/folder/%s/list", input.FolderID)).Send()
		if err != nil {
			return nil, err
		}

		defer rsp.Body().Close()

		if rsp.Status().IsError() {
			return nil, errors.New(rsp.Status().Text())
		}

		defer rsp.Body().Close()
		byts, err := io.ReadAll(rsp.Body().Raw())
		if err != nil {
			return nil, err
		}

		var body ListsResponse

		err = json.Unmarshal(byts, &body)
		if err != nil {
			return nil, err
		}

		return ctx.Respond(body.Lists, len(body.Lists))
	}

	form.SelectField("list_id", title).
		Placeholder("Select a list").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getLists)).
				WithFieldReference("folder-id", "folder-id").
				WithSearchSupport().
				WithPagination(10).
				End().
				RefreshOn("folder-id").
				GetDynamicSource(),
		).
		HelpText(desc)
}

func RegisterTasksInput(form *smartform.FormBuilder, title string, desc string, required bool) {
	getTasks := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		input := sdk.DynamicInputToType[struct {
			ListID string `json:"list-id,omitempty"`
		}](ctx)

		client := fastshot.NewClient(BaseURL).
			Auth().BearerToken(ctx.Auth().AccessToken).
			Header().
			AddAccept("application/json").
			Build()
		rsp, err := client.GET(fmt.Sprintf("/v2/list/%s/task", input.ListID)).Send()
		if err != nil {
			return nil, err
		}

		defer rsp.Body().Close()

		if rsp.Status().IsError() {
			return nil, errors.New(rsp.Status().Text())
		}

		defer rsp.Body().Close()
		byts, err := io.ReadAll(rsp.Body().Raw())
		if err != nil {
			return nil, err
		}

		var body TaskResponse

		err = json.Unmarshal(byts, &body)
		if err != nil {
			return nil, err
		}

		return ctx.Respond(body.Tasks, len(body.Tasks))
	}

	form.SelectField("task-id", title).
		Placeholder("Choose a task").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getTasks)).
				WithFieldReference("list-id", "list-id").
				WithSearchSupport().
				WithPagination(10).
				End().
				RefreshOn("list-id").
				GetDynamicSource(),
		).
		HelpText(desc)
}

func GetAssigneeInput(form *smartform.FormBuilder, title string, desc string, required bool) {
	getAssignees := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		input := sdk.DynamicInputToType[struct {
			ListID string `json:"list-id,omitempty"`
		}](ctx)

		client := fastshot.NewClient(BaseURL).
			Auth().BearerToken(ctx.Auth().AccessToken).
			Header().
			AddAccept("application/json").
			Build()

		rsp, err := client.GET(fmt.Sprintf("/v2/list/%s/member", input.ListID)).Send()
		if err != nil {
			return nil, err
		}

		defer rsp.Body().Close()

		if rsp.Status().IsError() {
			return nil, errors.New(rsp.Status().Text())
		}

		defer rsp.Body().Close()
		byts, err := io.ReadAll(rsp.Body().Raw())
		if err != nil {
			return nil, err
		}

		var body MembersResponse

		err = json.Unmarshal(byts, &body)
		if err != nil {
			return nil, err
		}

		members := body.Members

		items := arrutil.Map[Member, map[string]any](members, func(input Member) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.ID,
				"name": input.Username,
			}, true
		})

		return ctx.Respond(items, len(items))
	}

	form.SelectField("assignee-id", title).
		Placeholder("Choose a task").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getAssignees)).
				WithFieldReference("list-id", "list-id").
				WithSearchSupport().
				WithPagination(10).
				End().
				RefreshOn("list-id").
				GetDynamicSource(),
		).
		HelpText(desc)
}

func CreateItem(accessToken, name, url string) (map[string]interface{}, error) {
	fullURL := BaseURL + url
	data := []byte(fmt.Sprintf(`{
		"name": "%s"
	}`, name))

	req, err := http.NewRequest(http.MethodPost, fullURL, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(res.StatusCode)
	fmt.Println(string(body))

	var response map[string]interface{}
	if errs := json.Unmarshal(body, &response); errs != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return response, fmt.Errorf("error: %v", response)
	}

	return response, nil
}

func GetSpace(accessToken string, spaceID string) (map[string]interface{}, error) {
	url := "https://api.clickup.com/api/v2/space/" + spaceID

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get space from ClickUp API")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var respData map[string]interface{}

	if err := json.Unmarshal(body, &respData); err != nil {
		return nil, err
	}

	return respData, nil
}

var ClickupPriorityType = []*smartform.Option{
	{Value: "1", Label: "Urgent"},
	{Value: "2", Label: "High"},
	{Value: "3", Label: "Normal"},
	{Value: "4", Label: "Low"},
}

var ClickupOrderbyType = []*smartform.Option{
	{Value: "id", Label: "Id"},
	{Value: "created", Label: "Created"},
	{Value: "updated", Label: "Updated"},
	{Value: "due_date", Label: "Due Date"},
	{Value: "start_date", Label: "Start Date"},
}

func GetClickUpClient(accessToken string, endpoint string, method string, body interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s%s", BaseURL, endpoint)

	var req *http.Request
	var err error

	if body != nil && (method == http.MethodGet || method == http.MethodPut || method == http.MethodPatch) {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
	}

	req.Header.Set("Authorization", accessToken)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ClickUp API error (status code %d): %s", resp.StatusCode, string(bodyBytes))
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, err
	}

	return result, nil
}
