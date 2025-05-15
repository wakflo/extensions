package shared

import (
	"encoding/json"
	"errors"
	"io"
	"regexp"
	"time"

	"github.com/gookit/goutil/arrutil"
	"github.com/juicycleff/smartform/v1"
	fastshot "github.com/opus-domini/fast-shot"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"

	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

var (
	calendarForm = smartform.NewAuthForm("google-calendar-auth", "Google Calendar OAuth", smartform.AuthStrategyOAuth2)
	_            = calendarForm.
			OAuthField("oauth", "Google Calendar OAuth").
			AuthorizationURL("https://accounts.google.com/o/oauth2/auth").
			TokenURL("https://oauth2.googleapis.com/token").
			Scopes([]string{
			"https://www.googleapis.com/auth/calendar",
		}).
		Build()
)

var SharedGoogleCalendarAuth = calendarForm.Build()

// Function to clean the time string and return RFC3339 formatted time
func FormatTimeString(input string) string {
	// Regular expression to strip the timezone identifier (e.g., [Africa/Lagos])
	re := regexp.MustCompile(`\[[A-Za-z/_]+\]`)
	cleanedString := re.ReplaceAllString(input, "")

	// Parse the cleaned time string (which should now be RFC3339 compatible)
	parsedTime, err := time.Parse(time.RFC3339, cleanedString)
	if err != nil {
		// Return empty string or any default value in case of error
		return "invalid date format"
	}

	// Return the time formatted back to RFC3339
	return parsedTime.Format(time.RFC3339)
}

func RegisterCalendarProps(form *smartform.FormBuilder) {
	getCalendarID := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		client := fastshot.NewClient("https://www.googleapis.com/calendar/v3").
			Auth().BearerToken(authCtx.AccessToken).
			Header().
			AddAccept("application/json").
			Build()

		rsp, err := client.GET("/users/me/calendarList").Send()
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

		var calendarList CalendarList
		err = json.Unmarshal(byts, &calendarList)
		if err != nil {
			return nil, err
		}

		items := arrutil.Map[Calendar, map[string]any](calendarList.Items, func(input Calendar) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.ID,
				"name": input.Summary,
			}, true
		})
		return ctx.Respond(items, len(items))
	}

	form.SelectField("calendar_id", "Calendar").
		Placeholder("Select a calendar").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getCalendarID)).
				WithSearchSupport().
				End().
				GetDynamicSource(),
		).
		HelpText("Select a Google Calendar")
}

func RegisterCalendarEventProps(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getEventIDs := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		input := sdk.DynamicInputToType[struct {
			CalendarID string `json:"calendar_id"`
		}](ctx)

		client := fastshot.NewClient("https://www.googleapis.com/calendar/v3/calendars").
			Auth().BearerToken(authCtx.AccessToken).
			Header().
			AddAccept("application/json").
			Build()

		rsp, err := client.GET("/" + input.CalendarID + "/events").Send()
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

		var eventList EventList
		err = json.Unmarshal(byts, &eventList)
		if err != nil {
			return nil, err
		}

		events := eventList.Items
		items := arrutil.Map[Event, map[string]any](events, func(input Event) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.ID,
				"name": input.Summary,
			}, true
		})
		return ctx.Respond(items, len(items))
	}

	return form.SelectField("event_id", "Calendar Event").
		Placeholder("Select an event").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getEventIDs)).
				WithFieldReference("calendar_id", "calendar_id").
				WithSearchSupport().
				End().
				RefreshOn("calendar_id").
				GetDynamicSource(),
		).
		HelpText("Select a calendar event")
}
