package actions

import (
	"errors"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/pinterest/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getPinAnalyticsActionProps struct {
	BoardID     string `json:"board_id"`
	PinID       string `json:"pin_id"`
	MetricType  string `json:"metric_type"` // Changed from MetricTypes []string
	StartDate   string `json:"start_date,omitempty"`
	EndDate     string `json:"end_date,omitempty"`
	AppTypes    string `json:"app_types,omitempty"`
	SplitField  string `json:"split_field,omitempty"`
	AdAccountID string `json:"ad_account_id"`
}

type GetPinAnalyticsAction struct{}

// Metadata returns metadata about the action
func (a *GetPinAnalyticsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_pin_analytics",
		DisplayName:   "Get Pin Analytics",
		Description:   "Get analytics data for a specific pin",
		Type:          core.ActionTypeAction,
		Documentation: getPinAnalyticsDocs,
		Icon:          "mdi:chart-line",
		SampleOutput:  "",
		Settings:      core.ActionSettings{},
	}
}

func (a *GetPinAnalyticsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_pin_analytics", "Get Pin Analytics")
	form.Description("Retrieve analytics data for a specific pin including impressions, saves, clicks, and more")

	// Board selection
	shared.RegisterBoardsProps(form)

	// Pin selection (dependent on board)
	shared.RegisterBoardPinsProps(form)

	// Analytics Configuration Section
	form.SectionField("analytics_config", "Analytics Configuration")

	// Metric type - single select field
	form.SelectField("metric_type", "Metric Type").
		Required(true).
		AddOption("IMPRESSION", "Impressions").
		AddOption("SAVE", "Saves").
		AddOption("PIN_CLICK", "Pin Clicks").
		AddOption("OUTBOUND_CLICK", "Outbound Clicks").
		AddOption("TOTAL_COMMENTS", "Total Comments").
		AddOption("TOTAL_REACTIONS", "Total Reactions").
		AddOption("VIDEO_MRC_VIEW", "Video MRC Views").
		AddOption("VIDEO_AVG_WATCH_TIME", "Video Average Watch Time").
		AddOption("VIDEO_V50_WATCH_TIME", "Video V50 Watch Time").
		AddOption("QUARTILE_95_PERCENT_VIEW", "95% Quartile Views").
		AddOption("VIDEO_10S_VIEW", "10-Second Video Views").
		AddOption("VIDEO_START", "Video Starts")

	// Date fields - keep as DateField for better UX
	form.DateField("start_date", "Start Date").
		Required(true)

	form.DateField("end_date", "End Date").
		Required(true)

	// App types selection
	form.SelectField("app_types", "App Types").
		Required(false).
		AddOption("ALL", "All").
		AddOption("MOBILE", "Mobile").
		AddOption("TABLET", "Tablet").
		AddOption("WEB", "Web")

	// Split field for data segmentation
	form.SelectField("split_field", "Split Field").
		Required(false).
		AddOption("NO_SPLIT", "No Split").
		AddOption("APP_TYPE", "By App Type").
		AddOption("OWNED_PIN", "By Owned Pin").
		AddOption("SOURCE", "By Source").
		AddOption("AGE_BUCKET", "By Age Bucket").
		AddOption("GENDER", "By Gender")

	// Ad Account ID (optional)
	shared.GetAdAccountProps(form)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetPinAnalyticsAction) Auth() *core.AuthMetadata {
	return nil
}

// Helper function to format date to Pinterest's expected format (YYYY-MM-DD)
func formatDateForPinterest(dateStr string) (string, error) {
	// Try parsing common date formats
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02",
		"01/02/2006",
		"02-01-2006",
	}

	var parsedDate time.Time
	var err error

	for _, format := range formats {
		parsedDate, err = time.Parse(format, dateStr)
		if err == nil {
			break
		}
	}

	if err != nil {
		return "", errors.New("unable to parse date")
	}

	// Format to Pinterest's expected format
	return parsedDate.Format("2006-01-02"), nil
}

// Perform executes the action with the given context and input
func (a *GetPinAnalyticsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getPinAnalyticsActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.AccessToken == "" {
		return nil, errors.New("missing Pinterest auth token")
	}
	accessToken := authCtx.Token.AccessToken

	if input.BoardID == "" {
		return nil, errors.New("board ID is required")
	}

	if input.PinID == "" {
		return nil, errors.New("pin ID is required")
	}

	if input.MetricType == "" {
		return nil, errors.New("metric type is required")
	}

	// Validate and format dates
	if input.StartDate == "" {
		return nil, errors.New("start date is required")
	}

	if input.EndDate == "" {
		return nil, errors.New("end date is required")
	}

	// Convert dates to Pinterest format (YYYY-MM-DD)
	formattedStartDate, err := formatDateForPinterest(input.StartDate)
	if err != nil {
		return nil, errors.New("invalid start date format")
	}

	formattedEndDate, err := formatDateForPinterest(input.EndDate)
	if err != nil {
		return nil, errors.New("invalid end date format")
	}

	// Parse formatted dates for validation
	startDate, _ := time.Parse("2006-01-02", formattedStartDate)
	endDate, _ := time.Parse("2006-01-02", formattedEndDate)

	if startDate.After(endDate) {
		return nil, errors.New("start date must be before end date")
	}

	// Check if date range is not more than 90 days
	if endDate.Sub(startDate).Hours() > 90*24 {
		return nil, errors.New("date range cannot exceed 90 days")
	}

	// Validate ad account ID if provided
	if input.AdAccountID != "" {
		// Check if it contains only numbers
		for _, char := range input.AdAccountID {
			if char < '0' || char > '9' {
				return nil, errors.New("ad Account ID must contain only numbers")
			}
		}
	}

	// Convert single metric type to array for the API
	metricTypes := []string{input.MetricType}

	// Call the shared function to get pin analytics
	analytics, err := shared.GetPinAnalytics(
		accessToken,
		input.PinID,
		metricTypes, // Pass as array
		formattedStartDate,
		formattedEndDate,
		input.AppTypes,
		input.SplitField,
		input.AdAccountID,
	)
	if err != nil {
		return nil, err
	}

	return analytics, nil
}

func NewGetPinAnalyticsAction() sdk.Action {
	return &GetPinAnalyticsAction{}
}
