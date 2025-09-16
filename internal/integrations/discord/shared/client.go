package shared

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/juicycleff/smartform/v1"
	fastshot "github.com/opus-domini/fast-shot"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"

	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

const (
	APIBaseURL = "https://discord.com/api/v10"
)

func GetDiscordClient(token string, endpoint string, method string, body interface{}) ([]interface{}, error) {
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling body: %v", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	url := fmt.Sprintf("%s%s", APIBaseURL, endpoint)
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bot %s", token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "WakfloIntegration/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Discord API error: %s (Status Code: %d)", string(bodyBytes), resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if len(bodyBytes) == 0 {
		return []interface{}{}, nil
	}

	var arrayResult []interface{}
	if err := json.Unmarshal(bodyBytes, &arrayResult); err == nil {
		return arrayResult, nil
	}

	var objectResult map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &objectResult); err == nil {
		return []interface{}{objectResult}, nil
	}

	return nil, fmt.Errorf("error decoding response: could not parse as array or object")
}

func RegisterGuildsInput(form *smartform.FormBuilder, title string, desc string, required bool) {
	getGuilds := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}
		client := fastshot.NewClient("https://discord.com/api/v10").
			Auth().BearerToken("Bot " + authCtx.Extra["token"]).
			Header().
			AddAccept("application/json").
			Build()

		rsp, err := client.GET("/users/@me/guilds").Send()
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

		var guilds []Guild

		err = json.Unmarshal(byts, &guilds)
		if err != nil {
			return nil, err
		}

		return ctx.Respond(guilds, len(guilds))
	}

	form.SelectField("guild-id", title).
		Placeholder("Select a guild").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getGuilds)).
				WithSearchSupport().
				WithPagination(20).
				End().
				GetDynamicSource(),
		).
		HelpText(desc)
}

func RegisterChannelsInput(form *smartform.FormBuilder, title string, desc string, required bool) {
	getChannels := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		input := sdk.DynamicInputToType[struct {
			GuildID string `json:"guild-id,omitempty"`
		}](ctx)

		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}
		client := fastshot.NewClient("https://discord.com/api/v10").
			Auth().BearerToken("Bot " + authCtx.Extra["token"]).
			Header().
			AddAccept("application/json").
			Build()

		rsp, err := client.GET(fmt.Sprintf("/guilds/%s/channels", input.GuildID)).Send()
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

		var channels []Channel

		err = json.Unmarshal(byts, &channels)
		if err != nil {
			return nil, err
		}

		return ctx.Respond(channels, len(channels))
	}

	form.SelectField("channel-id", title).
		Placeholder("Select a channel").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getChannels)).
				WithFieldReference("guild-id", "guild-id").
				WithSearchSupport().
				WithPagination(10).
				End().
				RefreshOn("guild-id").
				GetDynamicSource(),
		).
		HelpText(desc)
}

func RegisterRolesInput(form *smartform.FormBuilder, title string, desc string, required bool) {
	getRoles := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		input := sdk.DynamicInputToType[struct {
			GuildID string `json:"guild-id,omitempty"`
		}](ctx)

		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}
		client := fastshot.NewClient("https://discord.com/api/v10").
			Auth().BearerToken("Bot " + authCtx.Extra["token"]).
			Header().
			AddAccept("application/json").
			Build()

		rsp, err := client.GET(fmt.Sprintf("/guilds/%s/roles", input.GuildID)).Send()
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

		var roles []Role // You'll need to define this struct

		err = json.Unmarshal(byts, &roles)
		if err != nil {
			return nil, err
		}

		return ctx.Respond(roles, len(roles))
	}

	form.SelectField("role-id", title).
		Placeholder("Select a role").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getRoles)).
				WithFieldReference("guild-id", "guild-id").
				WithSearchSupport().
				WithPagination(10).
				End().
				RefreshOn("guild-id").
				GetDynamicSource(),
		).
		HelpText(desc)
}

// CheckRateLimitResponse checks if we hit Discord's rate limit
func CheckRateLimitResponse(response map[string]interface{}) error {
	if response == nil {
		return nil
	}

	if message, ok := response["message"].(string); ok {
		if message == "You are being rate limited." {
			retryAfter, ok := response["retry_after"].(float64)
			if ok {
				return fmt.Errorf("rate limited, retry after %.2f seconds", retryAfter)
			}
			return errors.New("rate limited by Discord API")
		}
	}

	return nil
}
