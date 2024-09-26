package calendly

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	fastshot "github.com/opus-domini/fast-shot"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type User struct {
	URI                 string `json:"uri"`
	Name                string `json:"name"`
	Slug                string `json:"slug"`
	Email               string `json:"email"`
	SchedulingURL       string `json:"scheduling_url"`
	Timezone            string `json:"timezone"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
	CurrentOrganization string `json:"current_organization"`
	AvatarURL           string `json:"avatar_url"`
}
type UsersResponse struct {
	Users []User `json:"collection"`
}

type CurrentUserResponse struct {
	Resource User `json:"resource"`
}

var (
	// #nosec
	tokenURL   = "https://auth.calendly.com/oauth/token"
	sharedAuth = autoform.NewOAuthField("https://auth.calendly.com/oauth/authorize", &tokenURL, []string{}).Build()
)

const baseURL = "https://api.calendly.com"

// func getCalendlyUsersInput(title string, desc string, required bool) *sdkcore.AutoFormSchema {
// 	getUsers := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
// 		client := fastshot.NewClient(baseURL).
// 			Auth().BearerToken(ctx.Auth.AccessToken).
// 			Header().
// 			AddAccept("application/json").
// 			Build()

// 		rsp, err := client.GET("/users").Send()
// 		if err != nil {
// 			return nil, err
// 		}

// 		defer rsp.Body().Close()

// 		if rsp.Status().IsError() {
// 			return nil, errors.New(rsp.Status().Text())
// 		}

// 		byts, err := io.ReadAll(rsp.Body().Raw())
// 		if err != nil {
// 			return nil, err
// 		}

// 		var body UsersResponse

// 		err = json.Unmarshal(byts, &body)
// 		if err != nil {
// 			return nil, err
// 		}

// 		users := body.Users

// 		return arrutil.Map[User, map[string]any](users, func(input User) (target map[string]any, find bool) {
// 			return map[string]any{
// 				"id":   input.URI,
// 				"name": input.Name,
// 			}, true
// 		}), nil
// 	}

// 	return autoform.NewDynamicField(sdkcore.String).
// 		SetDisplayName(title).
// 		SetDescription(desc).
// 		SetDynamicOptions(&getUsers).
// 		SetRequired(required).Build()
// }

func getCurrentCalendlyUserInput(title string, desc string, required bool) *sdkcore.AutoFormSchema {
	getCurrentUser := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
		client := fastshot.NewClient(baseURL).
			Auth().BearerToken(ctx.Auth.AccessToken).
			Header().
			AddAccept("application/json").
			Build()

		rsp, err := client.GET("/users/me").Send()
		if err != nil {
			return nil, err
		}

		defer rsp.Body().Close()

		if rsp.Status().IsError() {
			return nil, errors.New(rsp.Status().Text())
		}

		byts, err := io.ReadAll(rsp.Body().Raw())
		if err != nil {
			return nil, err
		}

		var body CurrentUserResponse

		err = json.Unmarshal(byts, &body)
		if err != nil {
			return nil, err
		}

		user := body.Resource

		return []map[string]any{
			{
				"id":   user.URI,
				"name": user.Name,
			},
		}, nil
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(desc).
		SetDynamicOptions(&getCurrentUser).
		SetRequired(required).Build()
}

func listEvents(accessToken, url string, status string, user string) (map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("user", user)
	query.Add("status", status)

	req.URL.RawQuery = query.Encode()

	req.Header.Add("Authorization", "Bearer "+accessToken)
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

var calendlyEventStatusType = []*sdkcore.AutoFormSchema{
	{Const: "active", Title: "Active"},
	{Const: "canceled", Title: "Canceled"},
}
