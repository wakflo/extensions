package shared

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/gookit/goutil/arrutil"
	fastshot "github.com/opus-domini/fast-shot"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

var (
	// #nosec
	tokenURL   = "https://auth.calendly.com/oauth/token"
	SharedAuth = autoform.NewOAuthField("https://auth.calendly.com/oauth/authorize", &tokenURL, []string{}).Build()
)

const BaseURL = "https://api.calendly.com"

func GetCurrentCalendlyUserInput(title string, desc string, required bool) *sdkcore.AutoFormSchema {
	getCurrentUser := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		client := fastshot.NewClient(BaseURL).
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

		return ctx.Respond([]map[string]any{
			{
				"id":   user.URI,
				"name": user.Name,
			},
		}, 1)
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(desc).
		SetDynamicOptions(&getCurrentUser).
		SetRequired(required).Build()
}

func GetCalendlyEventInput(title string, desc string, required bool) *sdkcore.AutoFormSchema {
	getEvents := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		input := sdk.DynamicInputToType[struct {
			User string `json:"user"`
		}](ctx)

		client := fastshot.NewClient(BaseURL).
			Auth().BearerToken(ctx.Auth.AccessToken).
			Header().
			AddAccept("application/json").
			Build()

		req := client.GET("/scheduled_events")
		req.Query().AddParam("user", input.User)

		rsp, err := req.Send()
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

		var body EventsResponse

		err = json.Unmarshal(byts, &body)
		if err != nil {
			return nil, err
		}
		events := body.Events

		items := arrutil.Map[Event, map[string]any](events, func(input Event) (target map[string]any, find bool) {
			return map[string]any{
				"id":    input.URI,
				"name":  input.Name,
				"start": input.StartTime,
				"end":   input.EndTime,
			}, true
		})

		return ctx.Respond(items, len(items))
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(desc).
		SetDynamicOptions(&getEvents).
		SetRequired(required).
		Build()
}

func GetCalendlyEventTypeInput(title string, desc string, required bool) *sdkcore.AutoFormSchema {
	getEventTypes := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		input := sdk.DynamicInputToType[struct {
			User string `json:"user"`
		}](ctx)

		client := fastshot.NewClient(BaseURL).
			Auth().BearerToken(ctx.Auth.AccessToken).
			Header().
			AddAccept("application/json").
			Build()

		req := client.GET("/event_types")
		req.Query().AddParam("user", input.User)

		rsp, err := req.Send()
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

		var body EventTypesResponse

		err = json.Unmarshal(byts, &body)
		if err != nil {
			return nil, err
		}

		items := arrutil.Map[EventType, map[string]any](body.Collection, func(input EventType) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.URI,
				"name": input.Name,
			}, true
		})

		return ctx.Respond(items, len(items))
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(desc).
		SetDynamicOptions(&getEventTypes).
		SetRequired(required).
		Build()
}

func ListEvents(accessToken, url string, status string, user string) (map[string]interface{}, error) {
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

func GetEvent(accessToken, eventID string) (map[string]interface{}, error) {
	url := BaseURL + "/scheduled_events/" + getEventID(eventID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-OK status: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return response, nil
}

func CreateSingleUseLink(accessToken, eventTypeURI string, maxEventCount int) (map[string]interface{}, error) {
	url := BaseURL + "/scheduling_links"
	payload := map[string]interface{}{
		"max_event_count": maxEventCount,
		"owner":           eventTypeURI,
		"owner_type":      "EventType",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")

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

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return response, fmt.Errorf("error: %v", response)
	}

	return response, nil
}

var CalendlyEventStatusType = []*sdkcore.AutoFormSchema{
	{Const: "active", Title: "Active"},
	{Const: "canceled", Title: "Canceled"},
}

func getEventID(url string) string {
	// Define the regex pattern to capture the alphanumeric string after scheduled_events/
	re := regexp.MustCompile(`scheduled_events/([a-zA-Z0-9-]+)$`)
	match := re.FindStringSubmatch(url)

	// Check if there is a match
	if len(match) > 1 {
		return match[1]
	}
	return ""
}
