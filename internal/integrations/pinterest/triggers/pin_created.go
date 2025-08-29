package triggers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/pinterest/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type pinCreatedTriggerProps struct {
	BoardID  string `json:"board_id,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
}

type PinCreatedTrigger struct{}

func (t *PinCreatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "pin_created",
		DisplayName:   "Pin Created",
		Description:   "Triggers when a new pin is created in your Pinterest account",
		Type:          core.TriggerTypePolling,
		Documentation: pinCreatedDocs,
		SampleOutput: map[string]any{
			"id":               "123456789012345678",
			"created_at":       "2024-08-20T14:30:00Z",
			"link":             "https://www.pinterest.com/pin/123456789012345678",
			"title":            "Beautiful Sunset Photography",
			"description":      "Amazing sunset captured at the beach during golden hour",
			"alt_text":         "Sunset at beach",
			"note":             "Private note about this pin",
			"board_id":         "987654321098765432",
			"board_section_id": "876543210987654321",
			"media": map[string]any{
				"media_type": "image",
				"images": map[string]any{
					"150x150": map[string]any{
						"url":    "https://i.pinimg.com/150x150/example.jpg",
						"width":  150,
						"height": 150,
					},
					"400x300": map[string]any{
						"url":    "https://i.pinimg.com/400x300/example.jpg",
						"width":  400,
						"height": 300,
					},
					"600x": map[string]any{
						"url":    "https://i.pinimg.com/600x/example.jpg",
						"width":  600,
						"height": 450,
					},
				},
			},
			"parent_pin_id":     nil,
			"is_standard":       true,
			"has_been_promoted": false,
			"creative_type":     "REGULAR",
			"product_tags":      []any{},
		},
	}
}

func (t *PinCreatedTrigger) GetType() core.TriggerType {
	return core.TriggerTypePolling
}

func (t *PinCreatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("pin_created", "Pin Created")
	form.Description("Triggers when a new pin is created in your Pinterest account")

	// Configuration Section
	form.SectionField("configuration", "Configuration")

	// Use the shared board selector but make it optional
	boardField := shared.RegisterBoardsProps(form)
	boardField.Required(false).
		Placeholder("All Boards").
		HelpText("Leave empty to monitor all boards, or select a specific board to monitor")

	// Add page size configuration
	form.NumberField("page_size", "Page Size").
		Required(false).
		DefaultValue(25).
		HelpText("Number of pins to check per poll (max 100)")

	schema := form.Build()

	return schema
}

// Start initializes the PinCreatedTrigger, required for event and webhook triggers
func (t *PinCreatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Stop shuts down the PinCreatedTrigger
func (t *PinCreatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of PinCreatedTrigger
func (t *PinCreatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[pinCreatedTriggerProps](ctx)
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

	// Get last run time - handle the case where it doesn't exist
	var lastRunTime *time.Time
	lr, err := ctx.GetMetadata("lastRun")
	if err == nil && lr != nil {
		// Only type assert if we successfully got metadata
		if t, ok := lr.(*time.Time); ok {
			lastRunTime = t
		}
	}
	// If err != nil or type assertion fails, lastRunTime remains nil (first run)

	// Set default page size if not specified
	pageSize := input.PageSize
	if pageSize == 0 {
		pageSize = 25
	}

	var BASE_URL = "https://api.pinterest.com/v5"

	// Build the API URL
	var url string
	if input.BoardID != "" {
		// Get pins from specific board
		url = fmt.Sprintf("%s/boards/%s/pins?page_size=%d", BASE_URL, input.BoardID, pageSize)
	} else {
		// Get all user pins
		url = fmt.Sprintf("%s/pins?page_size=%d", BASE_URL, pageSize)
	}

	// Create a new HTTP GET request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Set required headers
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Initialize the HTTP client and execute the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Check for a successful response first
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("error fetching pins, status code: %d, body: %s", res.StatusCode, string(body))
	}

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response body
	var response struct {
		Items    []map[string]interface{} `json:"items"`
		Bookmark string                   `json:"bookmark,omitempty"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// If no last run time, store current time and return all pins (first run)
	if lastRunTime == nil {
		now := time.Now()
		if err := ctx.SetMetadata("lastRun", &now); err != nil {
			// Log error but don't fail the trigger
			fmt.Printf("Warning: failed to set lastRun metadata: %v\n", err)
		}
		return response.Items, nil
	}

	// Filter pins based on created_at time
	var filteredPins []map[string]interface{}
	var newestTime *time.Time

	for _, pin := range response.Items {
		createdAtStr, ok := pin["created_at"].(string)
		if !ok {
			continue // skip pins without created_at
		}

		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			continue // skip if timestamp is invalid
		}

		// Track the newest time for next run
		if newestTime == nil || createdAt.After(*newestTime) {
			newestTime = &createdAt
		}

		if createdAt.After(*lastRunTime) {
			filteredPins = append(filteredPins, pin)
		}
	}

	// Update lastRun metadata with the newest pin time
	if newestTime != nil {
		if err := ctx.SetMetadata("lastRun", newestTime); err != nil {
			fmt.Printf("Warning: failed to update lastRun metadata: %v\n", err)
		}
	}

	return filteredPins, nil
}

func (t *PinCreatedTrigger) Criteria(ctx context.Context) core.TriggerCriteria {
	return core.TriggerCriteria{}
}

func (t *PinCreatedTrigger) Auth() *core.AuthMetadata {
	return nil
}

func NewPinCreatedTrigger() sdk.Trigger {
	return &PinCreatedTrigger{}
}
