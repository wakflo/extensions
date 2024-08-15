package clickup

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var (
	// #nosec
	tokenURL   = "https://api.clickup.com/api/v2/oauth/token"
	sharedAuth = autoform.NewOAuthField("https://app.clickup.com/api", &tokenURL, []string{}).Build()
)

const baseURL = "https://api.clickup.com/api"

func getSpaces(accessToken, param string) (*Space, error) {
	reqURL := "https://api.clickup.com/api/v2/team/" + param + "/space"
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

	var respData Space
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return nil, err
	}
	var parsedResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	if err := json.Unmarshal(body, &parsedResponse); err != nil {
		return nil, err
	}

	return &Space{
		SpaceID: parsedResponse.ID,
		Name:    parsedResponse.Name,
	}, nil
}

func getData(accessToken, url string) (map[string]interface{}, error) {
	reqURL := url
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

func getList(accessToken, listID string) (map[string]interface{}, error) {
	reqURL := "https://api.clickup.com/api/v2/list/" + listID
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

func searchTask(accessToken, url string, page int, orderBy string, reverseOrder, includeClosed bool) (map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("page", strconv.Itoa(page))
	query.Add("order_by", orderBy)
	query.Add("reverse", strconv.FormatBool(reverseOrder))
	query.Add("include_closed", strconv.FormatBool(includeClosed))

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

type Team struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type Space struct {
	SpaceID string `json:"id"`
	Name    string `json:"name"`
}

func getTeams(accessToken string) ([]Team, error) {
	url := baseURL + "/v2/team"

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

	if err := json.Unmarshal(body, &parsedResponse); err != nil {
		return nil, err
	}

	return parsedResponse.Teams, nil
}

func createItem(accessToken, name, url string) (map[string]interface{}, error) {
	data := []byte(fmt.Sprintf(`{
		"name": "%s"
	}`, name))

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
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
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return response, fmt.Errorf("error: %v", response)
	}

	return response, nil
}

func getSpace(accessToken string, spaceID string) (map[string]interface{}, error) {
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

var clickupPriorityType = []*sdkcore.AutoFormSchema{
	{Const: "1", Title: "Urgent"},
	{Const: "2", Title: "High"},
	{Const: "3", Title: "Normal"},
	{Const: "4", Title: "Low"},
}

var clickupOrderbyType = []*sdkcore.AutoFormSchema{
	{Const: "id", Title: "Id"},
	{Const: "created", Title: "Created"},
	{Const: "updated", Title: "Updated"},
	{Const: "due_date", Title: "Due Date"},
	{Const: "start_date", Title: "Start Date"},
}
