package actions

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/juicycleff/smartform/v1"
	caption "github.com/lincaiyong/youtube-caption"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getVideoCaptionActionProps struct {
	VideoURL     string `json:"videoURL"`
	OutputFormat string `json:"outputFormat"`
	SaveSRT      bool   `json:"saveSRT"`
	SRTFilename  string `json:"srtFilename"`
}

type GetVideoCaptionAction struct{}

// Metadata returns metadata about the action
func (a *GetVideoCaptionAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_video_caption",
		DisplayName:   "Get YouTube Caption",
		Description:   "Extract captions/transcript from a YouTube video in whatever language is available. Supports multiple YouTube URL formats and SRT export.",
		Type:          core.ActionTypeAction,
		Documentation: getVideoCaptionDocs,
		Icon:          "youtube",
		SampleOutput: map[string]any{
			"video_id":    "dQw4w9WgXcQ",
			"video_url":   "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			"video_title": "Rick Astley - Never Gonna Give You Up",
			"language":    "en",
			"full_text":   "We're no strangers to love You know the rules and so do I...",
			"captions": []map[string]any{
				{
					"text":        "We're no strangers to love",
					"start_ms":    0,
					"duration_ms": 3000,
					"offset_text": "0:00",
				},
			},
			"srt_content":         "1\n00:00:00,000 --> 00:00:03,000\nWe're no strangers to love\n\n",
			"total_segments":      125,
			"available_languages": []string{"en", "es", "fr"},
			"srt_file_saved":      true,
			"srt_filename":        "Rick_Astley_-_Never_Gonna_Give_You_Up.srt",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GetVideoCaptionAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_video_caption", "Get YouTube Caption")

	form.TextField("videoURL", "YouTube Video URL").
		Required(true).
		HelpText("The YouTube video URL (e.g., https://www.youtube.com/watch?v=VIDEO_ID)")

	form.SelectField("outputFormat", "Output Format").
		Required(true).
		DefaultValue("both").
		AddOption("both", "Both (Structured + Full Text)").
		AddOption("structured", "Structured Captions Only").
		AddOption("text", "Full Text Only").
		AddOption("srt", "SRT Format").
		HelpText("Choose the format of the caption output")

	form.CheckboxField("saveSRT", "Save as SRT File").
		DefaultValue(false).
		HelpText("Save the captions as an SRT subtitle file")

	form.TextField("srtFilename", "SRT Filename").
		Required(false).
		HelpText("Filename for the SRT file (optional, will use video title if not specified)")

	schema := form.Build()

	return schema
}

// Auth returns nil since YouTube caption extraction doesn't require authentication
func (a *GetVideoCaptionAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action
func (a *GetVideoCaptionAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getVideoCaptionActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Extract video ID from URL
	videoID, err := extractVideoIDFallback(input.VideoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to extract video ID: %w", err)
	}

	// Get available caption tracks to show what languages are available
	tracks, err := caption.GetAvailableTracksWithContext(context.Background(), videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch caption tracks: %w", err)
	}

	if len(tracks) == 0 {
		return nil, fmt.Errorf("no captions available for this video")
	}

	// Get list of available languages
	availableLanguages := make([]string, 0, len(tracks))
	languageNames := make([]string, 0, len(tracks))
	for _, track := range tracks {
		availableLanguages = append(availableLanguages, track.LanguageCode)
		if track.Name.SimpleText != "" {
			languageNames = append(languageNames, track.Name.SimpleText)
		}
	}

	// Download captions using default options (library will pick the best available)
	opts := caption.DefaultOptions()
	captionData, err := caption.DownloadWithContext(context.Background(), videoID, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to download captions: %w", err)
	}

	// Get video title (placeholder for now)
	videoTitle := fmt.Sprintf("Video_%s", videoID)
	videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)

	// Try to determine which language was actually returned
	// Since the library doesn't tell us, we'll return the first available
	actualLanguage := "unknown"
	if len(availableLanguages) > 0 {
		actualLanguage = availableLanguages[0]
	}

	// Format the response
	result := map[string]interface{}{
		"video_id":            videoID,
		"video_url":           videoURL,
		"video_title":         videoTitle,
		"language":            actualLanguage,
		"total_segments":      len(captionData.Events),
		"available_languages": availableLanguages,
	}

	// Add language names if available
	if len(languageNames) > 0 {
		result["available_language_names"] = languageNames
	}

	// Generate SRT content using built-in method
	srtContent := captionData.GetSRT()

	// Format output based on selected format
	switch input.OutputFormat {
	case "text":
		result["full_text"] = captionData.GetPlainText()
	case "structured":
		result["captions"] = buildStructuredCaptionsFromEvents(captionData.Events)
	case "srt":
		result["srt_content"] = srtContent
	case "both":
		result["captions"] = buildStructuredCaptionsFromEvents(captionData.Events)
		result["full_text"] = captionData.GetPlainText()
	}

	// Save SRT file if requested
	if input.SaveSRT || input.OutputFormat == "srt" {
		filename := input.SRTFilename
		if filename == "" {
			// Use video title as filename, sanitizing it
			filename = sanitizeFilename(videoTitle)
		}

		err = saveSRTToFile(srtContent, filename)
		if err != nil {
			result["srt_save_error"] = err.Error()
		} else {
			result["srt_file_saved"] = true
			result["srt_filename"] = filename
			if !strings.HasSuffix(filename, ".srt") {
				result["srt_filename"] = filename + ".srt"
			}
		}
	}

	return result, nil
}

// buildStructuredCaptionsFromEvents converts caption.CaptionEvent to structured format
func buildStructuredCaptionsFromEvents(events []caption.CaptionEvent) []map[string]interface{} {
	structuredCaptions := make([]map[string]interface{}, 0)

	for _, event := range events {
		// Each event can have multiple segments
		if len(event.Segments) == 0 {
			continue
		}

		// Combine all segments in an event into one text
		var textBuilder strings.Builder
		for i, segment := range event.Segments {
			textBuilder.WriteString(segment.UTF8)
			if i < len(event.Segments)-1 {
				textBuilder.WriteString(" ")
			}
		}

		// Calculate duration (approximation - until next event or default duration)
		// Since we don't have direct duration info, we'll estimate
		durationMs := 3000 // Default 3 seconds

		structuredCaption := map[string]interface{}{
			"text":        textBuilder.String(),
			"start_ms":    event.TStartMs,
			"duration_ms": durationMs,
			"offset_text": formatTimestampFromMs(event.TStartMs),
		}
		structuredCaptions = append(structuredCaptions, structuredCaption)
	}

	return structuredCaptions
}

// formatTimestampFromMs formats milliseconds to MM:SS format
func formatTimestampFromMs(milliseconds int) string {
	totalSeconds := milliseconds / 1000
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// formatDuration formats time.Duration to human-readable string
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// extractVideoIDFallback is a fallback method to extract video ID
func extractVideoIDFallback(url string) (string, error) {
	url = strings.TrimSpace(url)

	// Handle different YouTube URL formats
	patterns := []struct {
		contains string
		split    string
	}{
		{"youtube.com/watch?v=", "v="},
		{"youtu.be/", "youtu.be/"},
		{"youtube.com/embed/", "embed/"},
		{"youtube.com/v/", "v/"},
		{"youtube.com/shorts/", "shorts/"},
	}

	for _, pattern := range patterns {
		if strings.Contains(url, pattern.contains) {
			parts := strings.Split(url, pattern.split)
			if len(parts) < 2 {
				continue
			}
			videoID := parts[1]

			// Clean up video ID - remove query params
			if idx := strings.IndexAny(videoID, "&?#"); idx != -1 {
				videoID = videoID[:idx]
			}

			if len(videoID) > 0 {
				return videoID, nil
			}
		}
	}

	// Check if it's just a video ID (11 characters)
	if len(url) == 11 && !strings.Contains(url, "/") && !strings.Contains(url, ".") {
		return url, nil
	}

	return "", fmt.Errorf("unsupported YouTube URL format or invalid video ID")
}

// saveSRTToFile saves SRT content to a file
func saveSRTToFile(srtContent string, filename string) error {
	// Ensure filename has .srt extension
	if !strings.HasSuffix(filename, ".srt") {
		filename += ".srt"
	}

	// Write to file
	return os.WriteFile(filename, []byte(srtContent), 0644)
}

// NewGetVideoCaptionAction creates a new instance of the action
func NewGetVideoCaptionAction() sdk.Action {
	return &GetVideoCaptionAction{}
}

func sanitizeFilename(filename string) string {
	// Remove or replace characters that are invalid in filenames
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
		"\n", " ",
		"\r", " ",
	)

	sanitized := replacer.Replace(filename)
	sanitized = strings.TrimSpace(sanitized)

	// Limit filename length
	if len(sanitized) > 200 {
		sanitized = sanitized[:200]
	}

	return sanitized
}
