package actions

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/juicycleff/smartform/v1"
	youtube "github.com/sh1nkey/youtube-downloader/v2"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getVideoCaptionActionProps struct {
	VideoURL     string `json:"videoURL"`
	OutputFormat string `json:"outputFormat"`
	Language     string `json:"language"`
	SaveSRT      bool   `json:"saveSRT"`
	SRTFilename  string `json:"srtFilename"`
}

type GetVideoCaptionAction struct{}

// Metadata returns metadata about the action
func (a *GetVideoCaptionAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_video_caption",
		DisplayName:   "Get YouTube Caption",
		Description:   "Extract captions/transcript from a YouTube video. Supports multiple YouTube URL formats and SRT export.",
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

	form.TextField("language", "Language Code").
		Required(false).
		DefaultValue("en").
		HelpText("Language code for captions (e.g., 'en' for English, 'es' for Spanish). Leave empty to use default.")

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
	videoID, err := youtube.ExtractVideoID(input.VideoURL)
	if err != nil {
		// If ExtractVideoID fails, try our custom extraction as fallback
		videoID, err = extractVideoIDFallback(input.VideoURL)
		if err != nil {
			return nil, fmt.Errorf("failed to extract video ID: %w", err)
		}
	}

	// Create YouTube client
	client := youtube.Client{
		HTTPClient: nil, // Uses default HTTP client
	}

	// Get video metadata
	video, err := client.GetVideoContext(context.Background(), input.VideoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch video: %w", err)
	}

	// Get list of available languages
	availableLanguages := make([]string, 0, len(video.CaptionTracks))
	for _, track := range video.CaptionTracks {
		availableLanguages = append(availableLanguages, track.LanguageCode)
	}

	// Determine language to use
	languageCode := input.Language
	if languageCode == "" {
		languageCode = "en"
	}

	// Get transcript for the specified language
	transcript, err := client.GetTranscriptCtx(context.Background(), video, languageCode)
	if err != nil {
		// If the exact language fails, try language prefix matching
		if err == youtube.ErrTranscriptDisabled {
			return nil, fmt.Errorf("transcripts are disabled for this video")
		}

		// Try to find a matching language
		foundLang := false
		for _, track := range video.CaptionTracks {
			if strings.HasPrefix(track.LanguageCode, languageCode) {
				transcript, err = client.GetTranscriptCtx(context.Background(), video, track.LanguageCode)
				if err == nil {
					languageCode = track.LanguageCode
					foundLang = true
					break
				}
			}
		}

		// If still not found and not English, try English as fallback
		if !foundLang && languageCode != "en" {
			for _, track := range video.CaptionTracks {
				if strings.HasPrefix(track.LanguageCode, "en") {
					transcript, err = client.GetTranscriptCtx(context.Background(), video, track.LanguageCode)
					if err == nil {
						languageCode = track.LanguageCode
						foundLang = true
						break
					}
				}
			}
		}

		if !foundLang {
			if len(availableLanguages) == 0 {
				return nil, fmt.Errorf("no captions available for this video")
			}
			return nil, fmt.Errorf("captions not available in language '%s'. Available languages: %s",
				input.Language, strings.Join(availableLanguages, ", "))
		}
	}

	// Format the response
	result := map[string]interface{}{
		"video_id":            videoID,
		"video_url":           input.VideoURL,
		"video_title":         video.Title,
		"video_author":        video.Author,
		"video_duration":      formatDuration(video.Duration),
		"language":            languageCode,
		"total_segments":      len(transcript),
		"available_languages": availableLanguages,
	}

	// Generate SRT content
	srtContent := convertToSRT(transcript)

	// Format output based on selected format
	switch input.OutputFormat {
	case "text":
		result["full_text"] = transcript.String()
	case "structured":
		result["captions"] = buildStructuredCaptions(transcript)
	case "srt":
		result["srt_content"] = srtContent
	case "both":
		result["captions"] = buildStructuredCaptions(transcript)
		result["full_text"] = transcript.String()
	}

	// Save SRT file if requested
	if input.SaveSRT || input.OutputFormat == "srt" {
		filename := input.SRTFilename
		if filename == "" {
			// Use video title as filename, sanitizing it
			filename = sanitizeFilename(video.Title)
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

// buildStructuredCaptions converts transcript to structured format
func buildStructuredCaptions(transcript youtube.VideoTranscript) []map[string]interface{} {
	captions := make([]map[string]interface{}, 0, len(transcript))

	for _, segment := range transcript {
		caption := map[string]interface{}{
			"text":        segment.Text,
			"start_ms":    segment.StartMs,
			"duration_ms": segment.Duration,
			"offset_text": segment.OffsetText,
		}
		captions = append(captions, caption)
	}

	return captions
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

// convertToSRT converts YouTube transcript to SRT format
func convertToSRT(transcript youtube.VideoTranscript) string {
	var srtBuffer bytes.Buffer

	for i, segment := range transcript {
		// SRT subtitle number (1-based)
		srtBuffer.WriteString(fmt.Sprintf("%d\n", i+1))

		// Convert timestamps to SRT format (HH:MM:SS,mmm --> HH:MM:SS,mmm)
		startTime := formatSRTTimestamp(segment.StartMs)
		endTime := formatSRTTimestamp(segment.StartMs + segment.Duration)
		srtBuffer.WriteString(fmt.Sprintf("%s --> %s\n", startTime, endTime))

		// Caption text
		srtBuffer.WriteString(segment.Text)
		srtBuffer.WriteString("\n\n") // Empty line between subtitles
	}

	return srtBuffer.String()
}

// formatSRTTimestamp converts milliseconds to SRT timestamp format (HH:MM:SS,mmm)
func formatSRTTimestamp(milliseconds int) string {
	totalSeconds := milliseconds / 1000
	ms := milliseconds % 1000

	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, seconds, ms)
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
