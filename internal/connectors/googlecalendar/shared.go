package googlecalendar

import (
	"encoding/json"
	"errors"
	"io"
	"regexp"
	"time"

	"github.com/gookit/goutil/arrutil"
	fastshot "github.com/opus-domini/fast-shot"
	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var (
	// #nosec
	tokenURL   = "https://oauth2.googleapis.com/token"
	sharedAuth = autoform.NewOAuthField("https://accounts.google.com/o/oauth2/auth", &tokenURL, []string{
		"https://www.googleapis.com/auth/calendar",
	}).Build()
)

// Function to clean the time string and return RFC3339 formatted time
func formatTimeString(input string) string {
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

func getCalendarInput(title string, desc string, required bool) *sdkcore.AutoFormSchema {
	getCalendarID := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
		client := fastshot.NewClient("https://www.googleapis.com/calendar/v3").
			Auth().BearerToken(ctx.Auth.AccessToken).
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

		return arrutil.Map[Calendar, map[string]any](calendarList.Items, func(input Calendar) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.ID,
				"name": input.Summary,
			}, true
		}), nil
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(desc).
		SetDynamicOptions(&getCalendarID).
		SetRequired(required).Build()
}

func getCalendarEventIDInput(title string, desc string, required bool) *sdkcore.AutoFormSchema {
	getEventIDs := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
		input := sdk.DynamicInputToType[struct {
			CalendarID string `json:"calendar_id"`
		}](ctx)

		client := fastshot.NewClient("https://www.googleapis.com/calendar/v3/calendars").
			Auth().BearerToken(ctx.Auth.AccessToken).
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
		return arrutil.Map[Event, map[string]any](events, func(input Event) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.ID,
				"name": input.Summary,
			}, true
		}), nil
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(desc).
		SetDynamicOptions(&getEventIDs).
		SetRequired(required).Build()
}
